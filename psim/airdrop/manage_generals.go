package airdrop

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
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

			ok, err := s.tryProcessGeneralAcc(ctx, acc)
			if err != nil {
				logger.WithError(err).Error("Failed to process GeneralAccount.")
				// Will try later
				s.pendingGeneralAccounts.Put(acc)
				break
			}

			if !ok {
				// This general Account is not ready for Issuance yet
				s.pendingGeneralAccounts.Put(acc)
				break
			}

			logger.WithField("account_address", acc).WithError(err).Info("Processed GeneralAccount successfully.")
		}
	}
}

func (s *Service) tryProcessGeneralAcc(ctx context.Context, accAddress string) (bool, error) {
	email, err := s.isReadyForIssuance(accAddress)
	if err != nil {
		if err == errUserNotFound {
			s.log.WithField("account_address", accAddress).
				Warn("Tried to check User's AirdropState, but User not found. I won't come back to this User again.")
			return true, nil
		}

		return false, errors.Wrap(err, "Failed to check readiness for issuance")
	}

	if email == "" {
		// Not ready for Issuance yet
		return false, nil
	}
	fields := logan.F{
		"email": email,
	}
	s.log.WithFields(fields).WithField("account_address", accAddress).Info("Found User, who is ready for Issuance.")

	// User is ready for Issuance
	balanceID, err := s.getBalanceID(accAddress)
	if err != nil {
		return false, errors.Wrap(err, "Failed to get BalanceID of the Account", fields)
	}
	fields["balance_id"] = balanceID

	err = s.processIssuance(ctx, accAddress, balanceID)
	if err != nil {
		return false, errors.Wrap(err, "Failed to process Issuance", fields)
	}

	err = s.sendEmail(email)
	if err != nil {
		// TODO Create separate routine, which will manage emails
		s.log.WithFields(fields).WithError(err).Error("Failed to send email.")
		// Don't return error, as the Issuance actually happened
	}

	return true, nil
}

var (
	errUserNotFound = errors.New("User not found.")
	errNoBalanceID  = errors.New("BalanceID not found for Account.")
)

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
