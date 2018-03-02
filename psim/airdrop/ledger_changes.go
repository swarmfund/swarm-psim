package airdrop

import (
	"time"

	"context"

	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/psim/psim/app"
)

func (s *Service) listenLedgerChangesInfinitely(ctx context.Context) {
	s.log.Info("Started listening Transactions stream.")
	txStream, txStreamerErrs := s.txStreamer.StreamTransactions(ctx)

	app.RunOverIncrementalTimer(ctx, s.log, "ledger_changes_processor", func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return nil
		case txEvent := <-txStream:
			if app.IsCanceled(ctx) {
				return nil
			}

			if txEvent.Transaction == nil {
				return nil
			}

			for _, change := range txEvent.Transaction.LedgerChanges() {
				s.processChange(ctx, txEvent.Transaction.CreatedAt, change)
			}

			return nil
		case txStreamerErr := <-txStreamerErrs:
			s.log.WithError(txStreamerErr).Error("TXStreamer sent error into its error channel.")
			return nil
		}
	}, 0, 10*time.Second)
}

func (s *Service) runOnce(ctx context.Context) {

}

func (s *Service) processChange(ctx context.Context, ts time.Time, change xdr.LedgerEntryChange) {
	switch change.Type {
	case xdr.LedgerEntryChangeTypeCreated:
		entryData := change.Created.Data

		if entryData.Type != xdr.LedgerEntryTypeAccount {
			return
		}

		accEntry := change.Created.Data.Account

		if ts.Sub(*s.config.RegisteredBefore) > 0 {
			// Account creation too late
			return
		}

		if accEntry.AccountType == xdr.AccountTypeGeneral {
			// Account was created already with General type
			s.streamGeneralAccount(ctx, accEntry.AccountId.Address())
			return
		} else {
			addr := accEntry.AccountId.Address()
			s.log.WithField("account_address", addr).Info("Found created Account.")
			s.createdAccounts[addr] = struct{}{}
			return
		}
	case xdr.LedgerEntryChangeTypeUpdated:
		entryData := change.Updated.Data

		if entryData.Type != xdr.LedgerEntryTypeAccount {
			return
		}

		accEntry := change.Updated.Data.Account

		if accEntry.AccountType != xdr.AccountTypeGeneral {
			// Account was updated but its Type is not General
			return
		}

		addr := accEntry.AccountId.Address()
		if _, ok := s.createdAccounts[addr]; ok {
			s.streamGeneralAccount(ctx, addr)
			delete(s.createdAccounts, addr)
		}
	}
}

func (s *Service) streamGeneralAccount(ctx context.Context, accAddress string) {
	select {
	case <-ctx.Done():
		return
	case s.generalAccountsCh <- accAddress:
		return
	}
}
