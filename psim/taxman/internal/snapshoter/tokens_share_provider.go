package snapshoter

import (
	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/psim/psim/taxman/internal/state"
)

//go:generate mockery -case underscore -testonly -inpkg -name tokenShareProvider
type tokenShareProvider interface {
	GetFeesToShareToTokenHolders() (map[state.AccountID]map[state.AssetCode]int64, error)
}

type tokenShareProviderImpl struct {
	state statable
}

func (p *tokenShareProviderImpl) GetFeesToShareToTokenHolders() (map[state.AccountID]map[state.AssetCode]int64, error) {
	totalTokensAmount := p.state.GetTotalTokensAmount()
	err := ensureAllNotNegative(totalTokensAmount)
	if err != nil {
		return nil, errors.Wrap(err, "invalid total tokens amount")
	}

	totalFeesToShare := p.state.GetTotalFeesToShare()
	err = ensureAllNotNegative(totalFeesToShare)
	if err != nil {
		return nil, errors.Wrap(err, "invalid total fees to share")
	}

	feesToShareToTokenHolders := map[state.AccountID]map[state.AssetCode]int64{}

	for tokenBalance := range p.state.TokenBalances() {
		totalTokens := totalTokensAmount[tokenBalance.Asset]
		if totalTokens == 0 {
			continue
		}

		if tokenBalance.Amount > totalTokens {
			return nil, fmt.Errorf("unexpected state amount %d exceeds total tokens amount %d", tokenBalance.Amount, totalTokens)
		}

		assetCode := p.state.GetAssetByToken(tokenBalance.Asset)
		amountOfFeesToShare := totalFeesToShare[assetCode]

		// We ignore overflow here on purpose, as it should never happened - tokenBalance.Amount/totalTokens is 1 at max.
		payoutAmount, _ := amount.BigDivide(amountOfFeesToShare, tokenBalance.Amount, totalTokens, amount.ROUND_DOWN)

		if _, ok := feesToShareToTokenHolders[tokenBalance.Account]; !ok {
			feesToShareToTokenHolders[tokenBalance.Account] = map[state.AssetCode]int64{}
		}

		feesToShareToTokenHolders[tokenBalance.Account][assetCode] += payoutAmount

	}

	return feesToShareToTokenHolders, nil
}

func ensureAllNotNegative(data map[state.AssetCode]int64) error {
	for asset, value := range data {
		if value < 0 {
			return fmt.Errorf("invalid amount %d for asset %s", value, asset)
		}
	}

	return nil
}
