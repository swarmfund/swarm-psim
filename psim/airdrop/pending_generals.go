package airdrop

import (
	"context"
	"time"

	"gitlab.com/swarmfund/psim/psim/app"
)

func (s *Service) processPendingGeneralAccounts(ctx context.Context) {
	pendingAccsTicker := time.Tick(10 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-pendingAccsTicker:
			if app.IsCanceled(ctx) {
				return
			}

			s.log.WithField("number_of_pending_accounts", s.pendingGeneralAccounts.Length()).Info("Started Processing pending GeneralAccounts batch.")

			var processedAccounts []string
			s.pendingGeneralAccounts.Range(ctx, func(acc string) {
				ok, err := s.tryProcessGeneralAcc(ctx, acc)
				if err != nil {
					s.log.WithField("account_address", acc).WithError(err).Error("Failed to process pending GeneralAccount.")
					return
				}

				if !ok {
					return
				}

				processedAccounts = append(processedAccounts, acc)
				s.log.WithField("account_address", acc).WithError(err).Info("Processed pending GeneralAccount successfully.")
				return
			})

			s.pendingGeneralAccounts.Delete(processedAccounts)
		}
	}
}
