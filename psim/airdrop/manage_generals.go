package airdrop

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
)

var (
	errUserNotFound = errors.New("User not found.")
	errNoBalanceID  = errors.New("BalanceID not found for Account.")
)

func (s *Service) consumeGeneralAccounts(ctx context.Context) {
	s.log.Info("Started consuming GeneralAccounts from stream.")

	for {
		select {
		case <-ctx.Done():
			return
		case acc := <-s.generalAccountsCh:
			if app.IsCanceled(ctx) {
				return
			}

			logger := s.log.WithFields(logan.F{
				"account_address": acc,
			})

			logger.Info("GeneralAccount was found.")

			email, err := s.isReadyForIssuance(acc)
			if err != nil {
				if err == errUserNotFound {
					s.log.WithField("account_address", acc).
						Warn("Tried to check User's AirdropState, but User not found. I won't come back to this User again.")
					break
				}

				logger.WithError(err).Error("Failed to check readiness for issuance.")
				// Will try later
				s.pendingGeneralAccounts.Put(acc)
				break
			}

			if email == "" {
				// Not ready for Issuance yet
				s.pendingGeneralAccounts.Put(acc)
				break
			}
			logger = logger.WithField("email", email)

			logger.Info("Found User, who is ready for Issuance.")

			err = s.processIssuance(ctx, acc, email)
			if err != nil {
				logger.WithError(err).Error("Failed to process Issuance for GeneralAccount. Will try later.")
				s.pendingGeneralAccounts.Put(acc)
				break
			}
		}
	}
}

func (s *Service) processIssuance(ctx context.Context, accAddress, email string) error {
	balanceID, err := s.getBalanceID(accAddress)
	if err != nil {
		return errors.Wrap(err, "Failed to get BalanceID of the Account")
	}
	fields := logan.F{"balance_id": balanceID}

	newIssuanceCreated, err := s.submitIssuance(ctx, accAddress, balanceID)
	if err != nil {
		return errors.Wrap(err, "Failed to process Issuance", fields)
	}

	if newIssuanceCreated {
		err = s.sendEmail(email)
		if err != nil {
			// TODO Create separate routine, which will manage emails
			s.log.WithFields(fields).WithError(err).Error("Failed to send email.")
			// Don't return error, as the Issuance actually happened
		}
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

func (s *Service) getBalanceID(accAddress string) (string, error) {
	balances, err := s.accountsConnector.Balances(accAddress)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get Account Balances")
	}

	for _, b := range balances {
		if b.Asset == s.config.Asset {
			return b.BalanceID, nil
		}
	}

	return "", errNoBalanceID
}

// TODO
func (s *Service) sendEmail(email string) error {
	return errors.New("Not implemented.")
}
