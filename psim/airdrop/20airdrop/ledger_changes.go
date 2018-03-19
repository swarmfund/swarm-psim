package mrefairdrop

import (
	"context"

	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/swarmfund/psim/psim/app"
)

func (s *Service) processChangesUpToSnapshotTime(ctx context.Context) {
	s.log.Info("Started listening TimedLedgers stream.")
	ledgerStream := s.ledgerStreamer.Run(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case timedLedger := <-ledgerStream:
			if app.IsCanceled(ctx) {
				return
			}

			if !s.processChange(ctx, timedLedger) {
				// Reached SnapshotTime - don't need to continue, whole job is done for this runner.
				s.log.WithField("snapshot_time", s.config.SnapshotTime).
					Info("Reached the SnapshotTime in the stream of LedgerEntryChanges.")
				return
			}
		}
	}
}

// ProcessChange only returns false if reached SnapshotTime.
func (s *Service) processChange(ctx context.Context, timedLedger airdrop.TimedLedgerChange) bool {
	if timedLedger.Time.Sub(s.config.SnapshotTime) > 0 {
		// Reached Snapshot time - Snapshot is fully ready all the following Changes, included this one are not interesting.
		return false
	}

	change := timedLedger.Change

	switch change.Type {
	case xdr.LedgerEntryChangeTypeCreated:
		s.processEntryData(change.Created.Data)
		return true
	case xdr.LedgerEntryChangeTypeUpdated:
		s.processEntryData(change.Updated.Data)
		return true
	default:
		// Not an Updated or Created type - not interested.
		return true
	}
}

func (s *Service) processEntryData(entryData xdr.LedgerEntryData) {
	if entryData.Type != xdr.LedgerEntryTypeBalance {
		return
	}

	entry := entryData.Balance

	if string(entry.Asset) != s.config.IssuanceAsset {
		// Not interested
		return
	}

	bonus, _ := s.snapshot[entry.AccountId.Address()]

	if bonus == nil {
		// This is probably master or comission Account or some other Account, which is created in genesis,
		// so no LedgerChange was found for this Account creation.
		return
	}

	bonus.BalanceID = entry.BalanceId.AsString()
	bonus.Balance = uint64(entry.Amount)
}
