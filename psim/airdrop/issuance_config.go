package airdrop

import (
	"unicode/utf8"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

type IssuanceConfig struct {
	Asset           string `fig:"asset,required"`
	Amount          uint64 `fig:"amount,required"`
	ReferenceSuffix string `fig:"reference_suffix"`
}

func (c IssuanceConfig) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"asset":            c.Asset,
		"amount":           c.Amount,
		"reference_suffix": c.ReferenceSuffix,
	}
}

func (c IssuanceConfig) Validate() (validationErr error) {
	if len(c.Asset) == 0 {
		return errors.New("Asset cannot be empty.")
	}

	if c.Amount == 0 {
		return errors.New("Amount cannot be zero.")
	}

	if utf8.RuneCountInString(c.ReferenceSuffix) > 8 {
		return errors.New("Reference suffix length must be at most 8 symbols.")
	}

	return nil
}
