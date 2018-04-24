package mrefairdrop

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/tokend/go/xdr"
)

// ProcessChangesUpToSnapshotTime is a blocking method, returns if ctx canceled or all the Changes are processed.
// Don't run this method in goroutine, as it won't notify anywhere when finished - will just return.
func (s *Service) processChangesUpToSnapshotTime(ctx context.Context) {
	s.log.WithField("snapshot_time", s.config.SnapshotTime).Info("Started listening TimedLedgers stream.")
	ledgerStream := s.ledgerStreamer.GetStream()
	// TODO Listen to Streamer stop
	go s.ledgerStreamer.Run(ctx)

	var snapshotPassAnnounced bool
	for {
		select {
		case <-ctx.Done():
			return
		case timedLedger := <-ledgerStream:
			if app.IsCanceled(ctx) {
				return
			}

			if timedLedger.Time.Sub(s.config.ApproveWaitFinishTime) > 0 {
				// Reached ApproveWaitFinishTime time - Snapshot is fully ready all the following Changes, included this one are not interesting.
				s.log.WithFields(logan.F{
					"approve_wait_finish_time": s.config.ApproveWaitFinishTime,
					"accounts_in_snapshot":     len(s.snapshot),
				}).Info("Reached the ApproveWaitFinishTime in the stream of LedgerEntryChanges.")
				return
			}

			if timedLedger.Time.Sub(s.config.SnapshotTime) > 0 {
				if !snapshotPassAnnounced {
					s.log.WithFields(logan.F{
						"snapshot_time":        s.config.SnapshotTime,
						"accounts_in_snapshot": len(s.snapshot),
					}).Info("Reached the SnapshotTime in the stream of LedgerEntryChanges.")
					snapshotPassAnnounced = true
				}

				s.processAfterSnapshotChange(ctx, timedLedger.Change)
				continue
			}

			// Before the Snapshot
			s.processChange(ctx, timedLedger.Change)
		}
	}
}

func (s *Service) processChange(ctx context.Context, change xdr.LedgerEntryChange) {
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
			return
		case xdr.LedgerEntryTypeBalance:
			// Balance created
			balEntry := entryData.Balance
			s.setBonusBalance(*balEntry)
			return
		default:
			return
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
					// Such a Referrer exists
					bonus.IsVerified = false
				}

				if accEntry.Referrer != nil {
					// Delete Referral as his AccountType is NotVerified now.
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

			return
		case xdr.LedgerEntryTypeBalance:
			balEntry := entryData.Balance
			s.setBonusBalance(*balEntry)
			return
		default:
			return
		}
	default:
		// Not an Updated or Created type - not interested.
		return
	}
}

func (s *Service) processAfterSnapshotChange(ctx context.Context, change xdr.LedgerEntryChange) {
	if change.Type != xdr.LedgerEntryChangeTypeUpdated {
		return
	}
	entryData := change.Updated.Data

	if entryData.Type != xdr.LedgerEntryTypeAccount {
		return
	}
	accEntry := entryData.Account

	switch accEntry.AccountType {
	case xdr.AccountTypeGeneral, xdr.AccountTypeSyndicate:
		// Account is probably becoming approved.
		bonus, ok := s.snapshot[accEntry.AccountId.Address()]
		if ok {
			// Such a Referrer exists.
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

	return
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
