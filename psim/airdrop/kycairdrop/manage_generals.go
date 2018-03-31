package kycairdrop

import (
	"context"

	"time"

	"fmt"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/issuance"
)

var (
	errUserNotFound = errors.New("User not found.")
)

func (s *Service) consumeGeneralAccounts(ctx context.Context) {
	s.log.Info("Started consuming GeneralAccounts from stream.")

	running.WithBackOff(ctx, s.log, "general_accounts_consumer", func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return nil
		case acc := <-s.generalAccountsCh:
			if app.IsCanceled(ctx) {
				return nil
			}

			logger := s.log.WithField("account_address", acc)

			if _, ok := s.blackList[acc]; ok {
				logger.Debug("Found GeneralAccount, but it's in BlackList - skipping it.")
				return nil
			}

			logger.Info("GeneralAccount was found, trying to process it.")

			running.UntilSuccess(ctx, logger, "general_account_processor", func(ctx context.Context) (bool, error) {
				err := s.processGeneralAccount(ctx, acc)
				if err != nil {
					return false, err
				} else {
					return true, nil
				}
			}, 5*time.Second, time.Hour)

			return nil
		}
	}, 0, 10*time.Second, time.Hour)
}

// ProcessGeneralAccount sends EmissionRequest and puts email into queue.
func (s *Service) processGeneralAccount(ctx context.Context, accAddress string) error {
	if s.isAlreadyIssued(accAddress) {
		s.log.WithField("account_id", accAddress).
			Info("Reference duplication - already processed Deposit, skipping.")
		return nil
	}

	emailAddress, err := s.getUserEmail(accAddress)
	if err != nil {
		if err == errUserNotFound {
			// Actually situation is not very probable.
			s.log.WithField("account_address", accAddress).
				Error("Tried to get User's emailAddress, but User not found. I won't come back to this User again.")
			// Returning nil, because we don't want to stop on this User and retry him.
			return nil
		}

		return errors.Wrap(err, "Failed to get User's emailAddress")
	}

	issuanceOpt, issuanceHappened, err := s.processIssuance(ctx, accAddress)
	if err != nil {
		return errors.Wrap(err, "Failed to process Issuance", logan.F{
			"email_address": emailAddress,
		})
	}

	logger := s.log.WithFields(logan.F{
		"account_address": accAddress,
		"email_address":   emailAddress,
		"issuance_opt":    *issuanceOpt,
	})
	if issuanceHappened {
		logger.Info("CoinEmissionRequest was sent successfully.")
	} else {
		logger.Info("Reference duplication in Horizon response - already processed Deposit, skipping.")
	}

	s.emailProcessor.AddEmailAddress(ctx, emailAddress)

	return nil
}

// TODO Return pointer to string and nil if no User existing. Avoid errUserNotFound.
func (s *Service) getUserEmail(accountAddress string) (email string, err error) {
	user, err := s.usersConnector.User(accountAddress)
	if err != nil {
		return "", errors.Wrap(err, "Failed to obtain User by accountAddress")
	}

	if user == nil {
		return "", errUserNotFound
	}

	return user.Attributes.Email, nil
}

func (s *Service) processIssuance(ctx context.Context, accAddress string) (*issuance.RequestOpt, bool, error) {
	balanceID, err := airdrop.GetBalanceID(accAddress, s.config.Asset, s.accountsConnector)
	if err != nil {
		return nil, false, errors.Wrap(err, "Failed to get BalanceID of the Account")
	}
	fields := logan.F{"balance_id": balanceID}

	opDetails := fmt.Sprintf(`{"cause": "%s"}`, airdrop.KYCIssuanceCause)
	issuanceOpt, ok, err := s.issuanceSubmitter.Submit(ctx, accAddress, balanceID, s.config.IssuanceConfig.Amount, opDetails)
	if err != nil {
		return nil, false, errors.Wrap(err, "Failed to process Issuance", fields)
	}

	return issuanceOpt, ok, nil
}
