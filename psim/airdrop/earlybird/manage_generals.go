package earlybird

import (
	"context"

	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/distributed_lab/running"
)

var (
	errUserNotFound = errors.New("User not found.")
)

func (s *Service) consumeGeneralAccounts(ctx context.Context) {
	s.log.Info("Started consuming GeneralAccounts from stream.")

	app.RunOverIncrementalTimer(ctx, s.log, "general_accounts_consumer", func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return nil
		case acc := <-s.generalAccountsCh:
			if running.IsCancelled(ctx) {
				return nil
			}

			s.log.WithField("account_address", acc).Info("GeneralAccount was found.")

			err := s.processGeneralAccount(ctx, acc)
			if err != nil {
				// Try this Account later
				s.pendingGeneralAccounts.Put(ctx, acc)
				return errors.Wrap(err, "Failed to process GeneralAccount, will try later")
			}

			return nil
		}
	}, 0, 10*time.Second)
}

func (s *Service) processGeneralAccount(ctx context.Context, accAddress string) error {
	logger := s.log.WithField("account_address", accAddress)

	email, err := s.isReadyForIssuance(accAddress)
	if err != nil {
		if err == errUserNotFound {
			s.log.WithField("account_address", accAddress).
				Warn("Tried to check User's AirdropState, but User not found. I won't come back to this User again.")
			return nil
		}

		return errors.Wrap(err, "Failed to check readiness for issuance.")
	}

	if email == "" {
		// Not ready for Issuance yet
		s.pendingGeneralAccounts.Put(ctx, accAddress)
		return nil
	}
	logger = logger.WithField("email", email)

	logger.Info("Found User, who is ready for Issuance.")

	err = s.processIssuance(ctx, accAddress, email)
	if err != nil {
		return errors.Wrap(err, "Failed to process Issuance for GeneralAccount. Will try later.", logan.F{
			"email": email,
		})
	}

	return nil
}

// IsReadyForIssuance returns empty email without error, if Account is not ready yet.
func (s *Service) isReadyForIssuance(accountAddress string) (email string, err error) {
	user, err := s.usersConnector.User(accountAddress)
	if err != nil {
		return "", errors.Wrap(err, "Failed to obtain User by accountAddress")
	}

	if user == nil {
		return "", errUserNotFound
	}

	if user.Attributes.AirdropState == "claimed" {
		return user.Attributes.Email, nil
	} else {
		// Not ready for Issuance yet
		return "", nil
	}
}

func (s *Service) processIssuance(ctx context.Context, accAddress, email string) error {
	balanceID, err := airdrop.GetBalanceID(accAddress, s.config.Asset, s.accountsConnector)
	if err != nil {
		return errors.Wrap(err, "Failed to get BalanceID of the Account")
	}
	fields := logan.F{"balance_id": balanceID}

	_, err = s.submitIssuance(ctx, accAddress, balanceID)
	if err != nil {
		return errors.Wrap(err, "Failed to process Issuance", fields)
	}

	s.emails.Put(ctx, email)

	return nil
}
