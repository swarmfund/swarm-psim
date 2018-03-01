package airdrop

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/app"
)

const (
	airdropStateClaimed = "claimed"
)

func (s *Service) processPendingGeneralAccounts(ctx context.Context) {
	s.log.Info("Started processing pending GeneralAccounts.")

	app.RunOverIncrementalTimer(ctx, s.log, "pending_generals_processor", func(ctx context.Context) error {
		if s.pendingGeneralAccounts.Length() == 0 {
			return nil
		}

		s.log.WithField("number_of_pending_accounts", s.pendingGeneralAccounts.Length()).Info("Started Processing pending GeneralAccounts batch.")

		processedAccounts := s.getAndProcessAPIUsers(ctx)
		if app.IsCanceled(ctx) {
			return nil
		}

		s.pendingGeneralAccounts.Delete(processedAccounts)
		return nil
	}, 10*time.Second, 10*time.Second)
}

// GetAndProcessAPIUsers does not return errors(only logged) intentionally, so that
// some single erroneous Issuance wouldn't prevent others to be processed.
func (s *Service) getAndProcessAPIUsers(ctx context.Context) (processedAccounts []string) {
	userStream, userStreamErrs := s.usersConnector.Users(ctx)

	for {
		select {
		case <-ctx.Done():
			return processedAccounts
		case user, ok := <-userStream:
			if !ok {
				// No more Users
				return processedAccounts
			}

			if user.Attributes.AirdropState != airdropStateClaimed {
				// Not ready for Issuance yet
				continue
			}

			if !s.pendingGeneralAccounts.Exists(user.ID) {
				// Not a pending Account
				continue
			}

			logger := s.log.WithFields(logan.F{
				"account_address": user.ID,
				"email":           user.Attributes.Email,
			})

			err := s.processIssuance(ctx, user.ID, user.Attributes.Email)
			if err != nil {
				logger.WithError(err).Error("Failed to process Issuance for pending GeneralAccount. Will try later.")
				continue
			}

			processedAccounts = append(processedAccounts, user.ID)
		case streamErr := <-userStreamErrs:
			s.log.WithError(streamErr).Error("Users connector sent error.")
		}
	}
}
