package btcdepositveri

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/externalsystems/derive"
	"gitlab.com/swarmfund/psim/psim/supervisor"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	DepositAsset       string `fig:"deposit_asset,required"`
	MinDepositAmount   uint64 `fig:"min_deposit_amount,required"`
	FixedDepositFee    uint64 `fig:"fixed_deposit_fee,required"`
	LastBlocksNotWatch uint64 `fig:"last_blocks_not_watch,required"`
	// NetworkType will be used to configure chain params
	NetworkType derive.NetworkType `fig:"network_type,required"`
	// ExternalSystem is optional and if provided will override value from deposit asset details
	ExternalSystem      int32  `fig:"external_system"`
	BlocksToSearchForTX uint64 `fig:"blocks_to_search_for_tx,required"`

	Source keypair.Address `fig:"source,required"`
	Signer keypair.Full    `fig:"signer,required"`
}

func NewConfig(configData map[string]interface{}) (*Config, error) {
	config := &Config{}

	err := figure.
		Out(config).
		From(configData).
		With(figure.BaseHooks, utils.CommonHooks, supervisor.DLFigureHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out")
	}

	return config, nil
}
