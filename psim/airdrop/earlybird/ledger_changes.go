package earlybird

import (
	"time"

	"context"

	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/internal"
	"gitlab.com/tokend/go/xdr"
)

func (s *Service) listenLedgerChangesInfinitely(ctx context.Context) {
	s.log.Info("Started listening Transactions stream.")
	txStream, txStreamerErrs := s.txStreamer.StreamTransactions(ctx)

	var isFirstTX = true
	var lastLoggedTXTime time.Time
	app.RunOverIncrementalTimer(ctx, s.log, "ledger_changes_processor", func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return nil
		case txEvent := <-txStream:
			if running.IsCancelled(ctx) {
				return nil
			}

			if txEvent.Transaction == nil {
				return nil
			}

			if isFirstTX {
				s.log.WithField("tx_time", txEvent.Meta.LatestLedger.ClosedAt).Info("Received first TX.")
				lastLoggedTXTime = txEvent.Meta.LatestLedger.ClosedAt
				isFirstTX = false
			} else {
				if txEvent.Meta.LatestLedger.ClosedAt.Sub(lastLoggedTXTime) > (24 * time.Hour) {
					// A day since last logged TX passed
					s.log.WithField("tx_time", txEvent.Meta.LatestLedger.ClosedAt).Info("Received next day TX.")
					lastLoggedTXTime = txEvent.Meta.LatestLedger.ClosedAt
				}
			}

			for _, change := range internal.LedgerChanges(txEvent.Transaction) {
				s.processChange(ctx, txEvent.Transaction.LedgerCloseTime, change)
			}

			return nil
		case txStreamerErr := <-txStreamerErrs:
			s.log.WithError(txStreamerErr).Error("TXStreamer sent error into its error channel.")
			return nil
		}
	}, 0, 10*time.Second)
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
