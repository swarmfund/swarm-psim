package snapshoter

import (
	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/psim/psim/taxman/internal/state"
)

//go:generate mockery -case underscore -testonly -inpkg -name payoutBuilder
type payoutBuilder interface {
	BuildOperations(result map[string]OperationSync, payoutInto map[state.AccountID]map[state.AssetCode]int64, typeOfPayout payoutType) error
}

type payoutBuilderImpl struct {
	state statable
}

func (b *payoutBuilderImpl) BuildOperations(result map[string]OperationSync,
	payoutInfo map[state.AccountID]map[state.AssetCode]int64, typeOfPayout payoutType) error {

	commissionAccount := b.state.GetSpecialAccounts().GetAccount(b.state.GetCommissionAccount())
	masterAccountID := b.state.GetMasterAccount()
	ledgerID := b.state.GetLedger()
	for accountID, amountToSharePerAsset := range payoutInfo {
		err := ensureAllNotNegative(amountToSharePerAsset)
		if err != nil {
			return errors.Wrap(err, "amount to share per asset has negative value")
		}

		for assetCode, amountToShare := range amountToSharePerAsset {
			if amountToShare == 0 {
				continue
			}

			account := b.state.GetAccount(accountID)
			balance := account.GetBalanceForAsset(masterAccountID, assetCode)
			if balance == nil {
				return fmt.Errorf("unexpected state: expected balance for %s %s %s to exist",
					accountID, masterAccountID, assetCode)
			}

			sourceBalance := commissionAccount.GetBalanceForAsset(masterAccountID, assetCode)
			if sourceBalance == nil {
				return fmt.Errorf("unexpected state: commission account balance for %s %s does not exist",
					masterAccountID, assetCode)
			}

			reference := reference(typeOfPayout, balance.Address, ledgerID)
			if _, exists := result[reference]; exists {
				return fmt.Errorf("unexpected state: reference already exists %s for %d %s %d",
					reference, typeOfPayout, balance.Address, ledgerID)
			}

			result[reference] = OperationSync{
				Reference:            reference,
				SourceBalanceID:      sourceBalance.Address,
				DestinationBalanceID: balance.Address,
				Amount:               amountToShare,
			}
		}
	}

	return nil
}
