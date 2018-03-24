package airdrop

import "gitlab.com/distributed_lab/logan/v3/errors"

type IssuanceConfig struct {
	Asset  string `fig:"asset,required"`
	Amount uint64 `fig:"amount,required"`
}

func (c IssuanceConfig) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"asset":  c.Asset,
		"amount": c.Amount,
	}
}

func (c IssuanceConfig) Validate() (validationErr error) {
	if len(c.Asset) == 0 {
		return errors.New("Asset cannot be empty.")
	}

	if c.Amount == 0 {
		return errors.New("Amount cannot be zero.")
	}

	return nil
}
