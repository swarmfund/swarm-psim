package airdrop

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
)

func (s *Service) consumeGeneralAccounts(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case acc := <-s.generalAccountsCh:
			if app.IsCanceled(ctx) {
				return
			}

			ok, err := s.tryProcessGeneralAcc(ctx, acc)
			if err != nil {
				s.log.WithField("account_address", acc).WithError(err).Error("Failed to process GeneralAccount.")
				// Will try later
				s.pendingGeneralAccounts.Put(acc)
				break
			}

			if !ok {
				// This general Account is not ready for Issuance yet
				s.pendingGeneralAccounts.Put(acc)
				break
			}

			s.log.WithField("account_address", acc).WithError(err).Info("Processed GeneralAccount successfully.")
		}
	}
}

func (s *Service) tryProcessGeneralAcc(ctx context.Context, accAddress string) (bool, error) {
	email, err := s.isReadyForIssuance(accAddress)
	if err != nil {
		return false, errors.Wrap(err, "Failed to check readiness for issuance")
	}

	if email == nil {
		return false, nil
	}
	fields := logan.F{
		"email": *email,
	}

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

	err = s.sendEmail(*email)
	if err != nil {
		// TODO Create separate routine, which will manage emails
		s.log.WithFields(fields).WithError(err).Error("Failed to send email.")
		// Don't return error, as the Issuance actually happened
	}

	return true, nil
}

// TODO
func (s *Service) isReadyForIssuance(userAddress string) (email *string, err error) {
	return nil, nil
}

// TODO
func (s *Service) getBalanceID(accAddress string) (string, error) {
	return "", nil
}

// TODO
func (s *Service) sendEmail(email string) error {
	return nil
}
