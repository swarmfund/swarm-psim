package btcdeposit

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/supervisor"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	Supervisor supervisor.Config `fig:"supervisor"`

	LastProcessedBlock uint64 `fig:"last_processed_block,required"`
	MinDepositAmount   uint64 `fig:"min_deposit_amount,required"`
	DepositAsset       string `fig:"deposit_asset,required"`
	OffchainCurrency   string `fig:"offchain_currency,required"`
	OffchainBlockchain string `fig:"offchain_blockchain,required"`
	// TODO Consider getting ExternalSystem integer by Asset code in runtime.
	ExternalSystem  int32  `fig:"external_system,required"`
	FixedDepositFee uint64 `fig:"fixed_deposit_fee,required"`
	DisableVerify   bool   `fig:"disable_verify"`

	Signer keypair.Full    `fig:"signer,required" mapstructure:"signer"`
	Source keypair.Address `fig:"source,required"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"supervisor": c.Supervisor,

		"last_processed_block":  c.LastProcessedBlock,
		"min_deposit_amount":    c.MinDepositAmount,
		"deposit_asset":         c.DepositAsset,
		"offchain_currency":     c.OffchainCurrency,
		"offchain_blockchain":   c.OffchainBlockchain,
		"external_system":       c.ExternalSystem,
		"fixed_deposit_fee":     c.FixedDepositFee,
		"disable_verify":        c.DisableVerify,
	}
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
