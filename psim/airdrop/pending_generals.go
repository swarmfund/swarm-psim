package airdrop

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
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

			if s.pendingGeneralAccounts.Length() == 0 {
				break
			}

			s.log.WithField("number_of_pending_accounts", s.pendingGeneralAccounts.Length()).Info("Started Processing pending GeneralAccounts batch.")

			processedAccounts := s.getAndProcessAPIUsers(ctx)
			if app.IsCanceled(ctx) {
				return
			}

			s.pendingGeneralAccounts.Delete(processedAccounts)
		}
	}
}

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

			if user.Attributes.AirdropState != "claimed" {
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

			logger.Info("Created IssuanceRequest for pending GeneralAccount successfully.")
			processedAccounts = append(processedAccounts, user.ID)
		case streamErr := <-userStreamErrs:
			s.log.WithError(streamErr).Error("Users connector sent error.")
		}
	}
}
