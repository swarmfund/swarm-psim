package snapshoter

import (
	"gitlab.com/tokend/psim/psim/taxman/internal/state"
	"gitlab.com/tokend/go/amount"
	"fmt"
)

//go:generate mockery -case underscore -testonly -inpkg -name referralShareProvider
type referralShareProvider interface {
	GetReferralSharePayout() (map[state.AccountID]map[state.AssetCode]int64, error)
}

type referralShareProviderImpl struct {
	state statable
}

func (p *referralShareProviderImpl) GetReferralSharePayout() (map[state.AccountID]map[state.AssetCode]int64, error) {
	result := map[state.AccountID]map[state.AssetCode]int64{}
	for child := range p.state.GetChildren() {

		if child.ShareForReferrer > amount.One * 100 || child.ShareForReferrer < 0 {
			return nil, fmt.Errorf("unexpected state: share for referrer for child %s is %d", child.Address, child.ShareForReferrer)
		}


		if _, ok := result[child.Parent]; !ok {
			result[child.Parent] = map[state.AssetCode]int64{}
		}

		for _, childBalance := range child.Balances {

			// overflow should never happen, as ShareForReferrer/100*amount.One is 1 at max
			feeToPayOut, _ := amount.BigDivide(
				childBalance.GetFeesToShare(), child.ShareForReferrer, amount.One*100, amount.ROUND_DOWN)
				feeToShareToParent := result[child.Parent][childBalance.Asset]

				if feeToShareToParent + feeToPayOut < feeToShareToParent {
					return nil, fmt.Errorf("overflow for %s", childBalance.Asset)
				}
			result[child.Parent][childBalance.Asset] += feeToPayOut
		}
	}

	return result, nil
}


