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

	for {
		select {
		case <-ctx.Done():
			return
		case txEvent := <-txStream:
			if app.IsCanceled(ctx) {
				return
			}

			if txEvent.Transaction == nil {
				break
			}

			for _, change := range txEvent.Transaction.LedgerChanges() {
				s.processChange(ctx, txEvent.Transaction.CreatedAt, change)
			}
		case txStreamerErr := <-txStreamerErrs:
			s.log.WithError(txStreamerErr).Error("TXStreamer sent error into its error channel.")
		}
	}
}

func (s *Service) processChange(ctx context.Context, ts time.Time, change xdr.LedgerEntryChange) {
	switch change.Type {
	case xdr.LedgerEntryChangeTypeCreated:
		entryData := change.Created.Data

		if entryData.Type == xdr.LedgerEntryTypeAccount {
			accEntry := change.Created.Data.Account

			if ts.Sub(s.config.RegisteredAfter) > 0 {
				// Account creation too late
				return
			}

			if accEntry.AccountType == xdr.AccountTypeGeneral {
				// Account was created already with General type
				s.streamGeneralAccount(ctx, accEntry.AccountId.Address())
			} else {
				s.createdAccounts[accEntry.AccountId.Address()] = struct{}{}
			}
		}
	case xdr.LedgerEntryChangeTypeUpdated:
		entryData := change.Updated.Data

		if entryData.Type == xdr.LedgerEntryTypeAccount {
			accEntry := change.Updated.Data.Account

			if accEntry.AccountType != xdr.AccountTypeGeneral {
				// Account was updated but its Type is not general
				return
			}

			addr := accEntry.AccountId.Address()
			if _, ok := s.createdAccounts[addr]; ok {
				s.streamGeneralAccount(ctx, addr)
				delete(s.createdAccounts, addr)
			}
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
