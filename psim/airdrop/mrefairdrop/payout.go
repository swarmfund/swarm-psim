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
	s.filterReferrals()
	s.filterReferrers()

	s.log.WithFields(logan.F{
		"accounts_in_snapshot": len(s.snapshot),
	}).Info("Started paying out airdrop according to to the Snapshot.")

	for accAddress, bonus := range s.snapshot {
		issAmount := countIssuanceAmount(len(bonus.Referrals), bonus.Balance)

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

			opDetails := fmt.Sprintf(`{"cause": "%s", "referrals": %d, "holdings": %d}`,
				airdrop.MarchReferralsIssuanceCause, len(bonus.Referrals), bonus.Balance)
			issuanceOpt, err := s.processIssuance(ctx, accAddress, bonus.BalanceID, issAmount, opDetails)
			if err != nil {
				return false, errors.Wrap(err, "Failed to process Issuance", logan.F{
					"email_address": *emailAddress,
				})
			}

			logger := s.log.WithFields(logan.F{
				"account_address": accAddress,
				"email_address":   *emailAddress,
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

func (s *Service) filterReferrers() {
	for accAddress, bonus := range s.snapshot {
		_, inBlackList := s.blackList[accAddress]

		if !bonus.IsVerified {
			delete(s.snapshot, accAddress)
			continue
		}

		if len(bonus.Referrals) == 0 {
			// There's nothing to payout airdrop for.
			delete(s.snapshot, accAddress)
			continue
		}

		if inBlackList {
			s.log.WithField("account_address", accAddress).Info("Filtering out Account, because it's in BlackList")
			delete(s.snapshot, accAddress)
			continue
		}
	}
}

func (s *Service) filterReferrals() {
	for referrer, bonus := range s.snapshot {
		for referral, _ := range bonus.Referrals {
			if _, inBlackList := s.blackList[referral]; inBlackList {
				s.log.WithFields(logan.F{
					"referrer_address": referrer,
					"referral_address": referral,
				}).Info("Filtering out referral, because it's in BlackList")
				delete(bonus.Referrals, referral)
			}
		}
	}
}

func countIssuanceAmount(referrals int, balance uint64) uint64 {
	result := balance * uint64(referrals) / 100

	if result > 4000*amount.One {
		result = 4000 * amount.One
	}

	result += 5 * amount.One * uint64(referrals)

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
