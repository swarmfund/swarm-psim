package kycairdrop

import (
	"context"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/lchanges"
	"sync"
)

// ListenLedgerChanges listens the stream from LedgerStreamer infinitely
// and processes each incoming LedgerChange.
// ListenLedgerChanges is a blocking method, it only returns when ctx is cancelled.
func (s *Service) listenLedgerChanges(ctx context.Context) {
	ledgerStream := s.ledgerStreamer.GetStream()
	lStreamerWG := sync.WaitGroup{}
	go func() {
		lStreamerWG.Add(1)
		s.ledgerStreamer.Run(ctx)
		lStreamerWG.Done()
	}()
	s.log.Info("Started listening TimedLedgers stream.")

	defer lStreamerWG.Wait()

	for {
		select {
		case <-ctx.Done():
			return
		case timedLedger := <-ledgerStream:
			if app.IsCanceled(ctx) {
				return
			}

			s.processChange(ctx, timedLedger)
		}
	}
}

func (s *Service) processChange(ctx context.Context, timedLedger lchanges.TimedLedgerChange) {
	change := timedLedger.Change

	switch change.Type {
	case xdr.LedgerEntryChangeTypeCreated:
		entryData := change.Created.Data

		if entryData.Type != xdr.LedgerEntryTypeAccount {
			return
		}

		accEntry := change.Created.Data.Account

		if accEntry.AccountType != xdr.AccountTypeGeneral {
			// Account of a non-General type was created - not interested.
			return
		}

		// Account was created already with General type.
		s.streamGeneralAccount(ctx, accEntry.AccountId.Address())
		return
	case xdr.LedgerEntryChangeTypeUpdated:
		entryData := change.Updated.Data

		if entryData.Type != xdr.LedgerEntryTypeAccount {
			return
		}

		accEntry := change.Updated.Data.Account

		if accEntry.AccountType != xdr.AccountTypeGeneral {
			// Account was updated but its Type is not General - not interested.
			return
		}

		addr := accEntry.AccountId.Address()
		s.streamGeneralAccount(ctx, addr)
		return
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
