package taxman

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/psim/psim/taxman/internal/state"
)

func (s *Service) SetUpListener() {
	for ; ; <-s.ticker.C {
		// getting commission and operational account balances
		info, err := s.horizon.Info()
		if err != nil {
			s.errors <- errors.Wrap(err, "failed to get horizon info")
			continue
		}

		// we do not consider operational account as special, as it participates in all the calculations as regular account
		// so even though operational account requires special handing by its nature it's closer to regular account
		operationalAccountID := state.AccountID(info.OperationalAccountID)
		s.state.SetOperationalAccount(operationalAccountID)
		if !s.state.Exists(operationalAccountID) {
			operationalAccount, err := s.createSpecialAccount(operationalAccountID)
			if err != nil {
				s.errors <- errors.Wrap(err, "failed to init operational account")
				continue
			}

			s.state.AddAccount(operationalAccount)
		}

		s.state.SetSpecialAccounts(state.AccountID(info.MasterAccountID),
			state.AccountID(info.StorageFeeAccountID), state.AccountID(info.CommissionAccountID))
		err = s.initStateWithSpecialAccounts(s.state.GetMasterAccount(), s.state.GetStorageFeeAccount(),
			s.state.GetCommissionAccount())
		if err != nil {
			s.errors <- errors.Wrap(err, "failed to init special accounts")
			continue
		}

		//_, err = s.GetPayoutCounter()
		//if err != nil {
		//	s.errors <- err
		//	continue
		//}
		return
	}
}

func (s *Service) initStateWithSpecialAccounts(specialAccounts ...state.AccountID) error {
	for _, specialAccountID := range specialAccounts {
		if s.state.GetSpecialAccounts().Exists(specialAccountID) {
			continue
		}

		specialAccount, err := s.createSpecialAccount(specialAccountID)
		if err != nil {
			return errors.Wrap(err, "failed to create special account")
		}

		s.state.GetSpecialAccounts().AddAccount(specialAccount)
	}

	return nil
}

// initSpecialAccount - creates fully valid instance of SpecialAccount with balances
func (s *Service) createSpecialAccount(specialAccountID state.AccountID) (state.Account, error) {
	accountInfo, err := s.horizon.AccountSigned(s.config.Signer, string(specialAccountID))
	if err != nil {
		return state.Account{}, errors.Wrap(err, "failed to get special account from horizon")
	}

	if accountInfo == nil {
		return state.Account{}, errors.New("special account not found")
	}

	account := state.Account{
		Address: specialAccountID,
	}
	for _, horizonBalance := range accountInfo.Balances {
		balance := state.Balance{
			Account:    account.Address,
			Address:    state.BalanceID(horizonBalance.BalanceID),
			Asset:      state.AssetCode(horizonBalance.Asset),
			ExchangeID: state.AccountID(horizonBalance.ExchangeID),
		}
		err := account.AddBalance(balance)
		if err != nil {
			return state.Account{}, errors.Wrap(err, "failed to add balance")
		}
	}

	return account, nil
}

func (s *Service) Listener(ctx context.Context) {
	s.SetUpListener()

	events := s.sse.Events()
	for {
		select {
		case event := <-events:
			if event.Err != nil {
				s.errors <- errors.Wrap(event.Err, "failed to get event")
				continue
			}

			critical, err := s.processEvent(event)
			if err != nil {
				s.errors <- errors.Wrap(err, "failed to process event")
				if critical {
					s.TearDown()
					return
				}
				continue
			}
			break
		case <-s.ctx.Done():
			s.log.Info("shutting down listener")
			// TODO sse does not allow to clean up connection
			return
		}
	}
}
