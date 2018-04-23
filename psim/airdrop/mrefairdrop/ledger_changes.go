package mrefairdrop

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/swarmfund/psim/psim/app"
)

// ProcessChangesUpToSnapshotTime is a blocking method, returns if ctx canceled or all the Changes are processed.
// Don't run this method in goroutine, as it won't notify anywhere when finished - will just return.
func (s *Service) processChangesUpToSnapshotTime(ctx context.Context) {
	s.log.WithField("snapshot_time", s.config.SnapshotTime).Info("Started listening TimedLedgers stream.")
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
				s.log.WithFields(logan.F{
					"snapshot_time":        s.config.SnapshotTime,
					"accounts_in_snapshot": len(s.snapshot),
				}).Info("Reached the SnapshotTime in the stream of LedgerEntryChanges.")
				return
			}
		}
	}
}

// TODO Stop returning bool, check for time above; pass LedgerChange without time.
// ProcessChange only returns false if reached SnapshotTime.
func (s *Service) processChange(ctx context.Context, timedLedger airdrop.TimedLedgerChange) bool {
	if timedLedger.Time.Sub(s.config.SnapshotTime) > 0 {
		// Reached Snapshot time - Snapshot is fully ready all the following Changes, included this one are not interesting.
		return false
	}

	change := timedLedger.Change

	switch change.Type {
	case xdr.LedgerEntryChangeTypeCreated:
		entryData := change.Created.Data

		switch entryData.Type {
		case xdr.LedgerEntryTypeAccount:
			// Account created
			accEntry := entryData.Account

			bonus := newBonusParams()
			if accEntry.AccountType == xdr.AccountTypeGeneral || accEntry.AccountType == xdr.AccountTypeSyndicate {
				// Account is created in already approved type.
				bonus.IsVerified = true

				if accEntry.Referrer != nil {
					referrerBonus, ok := s.snapshot[accEntry.Referrer.Address()]
					if ok {
						referrerBonus.addReferral(accEntry.AccountId.Address())
					}
				}
			}

			s.snapshot[accEntry.AccountId.Address()] = &bonus
			return true
		case xdr.LedgerEntryTypeBalance:
			// Balance created
			balEntry := entryData.Balance
			s.setBonusBalance(*balEntry)
			return true
		default:
			return true
		}
	case xdr.LedgerEntryChangeTypeUpdated:
		entryData := change.Updated.Data

		switch entryData.Type {
		case xdr.LedgerEntryTypeAccount:
			accEntry := entryData.Account

			switch accEntry.AccountType {
			case xdr.AccountTypeNotVerified:
				// Account could become not approved.
				bonus, ok := s.snapshot[accEntry.AccountId.Address()]
				if ok {
					bonus.IsVerified = false
				}

				if accEntry.Referrer != nil {
					// Delete Referral as his AccountType is NotVerified.
					referrerBonus, ok := s.snapshot[accEntry.Referrer.Address()]
					if ok {
						referrerBonus.deleteReferral(accEntry.AccountId.Address())
					}
				}
			case xdr.AccountTypeGeneral, xdr.AccountTypeSyndicate:
				// Account is probably becoming approved.
				bonus, ok := s.snapshot[accEntry.AccountId.Address()]
				if ok {
					bonus.IsVerified = true
				}

				if accEntry.Referrer != nil {
					// Add Referral as he became approved..
					referrerBonus, ok := s.snapshot[accEntry.Referrer.Address()]
					if ok {
						referrerBonus.addReferral(accEntry.AccountId.Address())
					}
				}
			}

			return true
		case xdr.LedgerEntryTypeBalance:
			balEntry := entryData.Balance
			s.setBonusBalance(*balEntry)
			return true
		default:
			return true
		}
	default:
		// Not an Updated or Created type - not interested.
		return true
	}
}

func (s *Service) setBonusBalance(entry xdr.BalanceEntry) {
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
