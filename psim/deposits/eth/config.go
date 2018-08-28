package eth

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/keypair"
)

type DepositConfig struct {
	Source              keypair.Address `fig:"source,required"`
	Signer              keypair.Full    `fig:"signer,required"`
	Cursor              uint64          `fig:"cursor,required"`
	ExternalSystem      int32           `fig:"external_system"`
	DepositAsset        string          `fig:"deposit_asset,required"`
	MinDepositAmount    uint64          `fig:"min_deposit_amount,required"`
	FixedDepositFee     uint64          `fig:"fixed_deposit_fee,required"`
	BlocksToSearchForTX uint64          `fig:"blocks_to_search_for_tx,required"`
}

func NewDepositConfig(raw map[string]interface{}) (*DepositConfig, error) {
	var config DepositConfig
	err := figure.
		Out(&config).
		With(figure.BaseHooks, utils.ETHHooks).
		From(raw).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "failed to figure out")
	}
	return &config, nil
}

type DepositVerifyConfig struct {
	Source              keypair.Address `fig:"source,required"`
	Signer              keypair.Full    `fig:"signer,required"`
	Cursor              uint64          `fig:"cursor,required"`
	ExternalSystem      int32           `fig:"external_system"`
	DepositAsset        string          `fig:"deposit_asset,required"`
	MinDepositAmount    uint64          `fig:"min_deposit_amount,required"`
	FixedDepositFee     uint64          `fig:"fixed_deposit_fee,required"`
	BlocksToSearchForTX uint64          `fig:"blocks_to_search_for_tx,required"`
	Confirmations       uint64          `fig:"confirmations,required"`
}

func NewDepositVerifyConfig(raw map[string]interface{}) (*DepositVerifyConfig, error) {
	var config DepositVerifyConfig
	err := figure.
		Out(&config).
		With(figure.BaseHooks, utils.ETHHooks).
		From(raw).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "failed to figure out")
	}
	return &config, nil
}
