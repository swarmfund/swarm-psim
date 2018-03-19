package mrefairdrop

import (
	"context"
	"time"

	"fmt"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/swarmfund/psim/psim/issuance"
)

func (s *Service) payOutSnapshot(ctx context.Context) {
	s.log.Info("Started paying out airdrop according to to the Snapshot.")

	s.filterBadAccounts()

	for accAddress, bonus := range s.snapshot {
		issAmount := countIssuanceAmount(bonus.Balance)

		if issAmount == 0 {
			continue
		}

		running.UntilSuccess(ctx, s.log, "general_account_processor", func(ctx context.Context) (bool, error) {
			emailAddress, err := s.getUserEmail(accAddress)
			if err != nil {
				return false, errors.Wrap(err, "Failed to get User's emailAddress")
			}
			if emailAddress == nil {
				// Actually situation is not very probable.
				s.log.WithField("account_address", accAddress).
					Error("Tried to get User's emailAddress, but User not found. I won't come back to this User again.")
				// Returning nil, because we don't want to stop on this User and retry him.
				return true, nil
			}

			opDetails := fmt.Sprintf(`{"cause": "%s", "holdings": %d}`, airdrop.March20b20IssuanceCause, bonus.Balance)
			issuanceOpt, err := s.processIssuance(ctx, accAddress, bonus.BalanceID, issAmount, opDetails)
			if err != nil {
				return false, errors.Wrap(err, "Failed to process Issuance", logan.F{
					"email_address": emailAddress,
				})
			}

			logger := s.log.WithFields(logan.F{
				"account_address": accAddress,
				"email_address":   emailAddress,
			})
			if issuanceOpt != nil {
				logger.WithField("issuance_opt", *issuanceOpt).Info("CoinEmissionRequest was sent successfully.")
			} else {
				logger.Info("Reference duplication - already processed Deposit, skipping.")
			}

			s.emailProcessor.AddEmailAddress(ctx, *emailAddress)

			return true, nil
		}, 5*time.Second, 5*time.Minute)
	}
}

func (s *Service) filterBadAccounts() {
	for accAddress, bonus := range s.snapshot {
		_, inBlackList := s.blackList[accAddress]

		if inBlackList || !bonus.IsVerified {
			delete(s.snapshot, accAddress)
		}
	}
}

func countIssuanceAmount(balance uint64) uint64 {
	result := balance * 20 / 100

	if result > 20000*amount.One {
		result = 20000 * amount.One
	}

	return result
}

func (s *Service) getUserEmail(accountAddress string) (email *string, err error) {
	user, err := s.usersConnector.User(accountAddress)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to obtain User by accountAddress")
	}

	if user == nil {
		return nil, nil
	}

	return &user.Attributes.Email, nil
}

func (s *Service) processIssuance(ctx context.Context, accAddress, balanceID string, amount uint64, opDetails string) (*issuance.RequestOpt, error) {
	var err error

	if balanceID == "" {
		s.log.WithFields(logan.F{
			"account_address": accAddress,
			"issuance_amount": amount,
		}).Warn("Found Issuance to be processed without BalanceID provided, we are ready for this, but it shouldn't have happened.")

		bIdP, err := s.balanceIDProvider.GetBalanceID(accAddress, s.config.IssuanceAsset)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get BalanceID of the Account")
		}
		if bIdP == nil {
			return nil, errors.From(errors.New("Account does not have a Balance for the provided asset."), logan.F{
				"asset": s.config.IssuanceAsset,
			})
		}
	}

	fields := logan.F{"balance_id": balanceID}

	issuanceOpt, err := s.issuanceSubmitter.Submit(ctx, accAddress, balanceID, amount, opDetails)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to process Issuance", fields)
	}

	return issuanceOpt, nil
}
