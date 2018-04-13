package btcdepositveri

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	Host string `fig:"host,required"`
	Port int    `fig:"port,required"`

	DepositAsset       string `fig:"deposit_asset,required"`
	MinDepositAmount   uint64 `fig:"min_deposit_amount,required"`
	FixedDepositFee    uint64 `fig:"fixed_deposit_fee,required"`
	LastBlocksNotWatch uint64 `fig:"last_blocks_not_watch,required"`

	Signer keypair.Full `fig:"signer,required" mapstructure:"signer"`
}

func NewConfig(configData map[string]interface{}) (*Config, error) {
	config := &Config{}

	err := figure.
		Out(config).
		From(configData).
		With(figure.BaseHooks, utils.CommonHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out")
	}

	return config, nil
}
