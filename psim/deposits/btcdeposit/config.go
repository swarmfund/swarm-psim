package btcdeposit

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/externalsystems/derive"
	"gitlab.com/swarmfund/psim/psim/supervisor"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/keypair"
)

// FIXME implement validate
type Config struct {
	Supervisor supervisor.Config `fig:"supervisor"`

	LastProcessedBlock uint64 `fig:"last_processed_block,required"`
	MinDepositAmount   uint64 `fig:"min_deposit_amount,required"`
	DepositAsset       string `fig:"deposit_asset,required"`
	// NetworkType will be used to configure chain params
	NetworkType derive.NetworkType `fig:"network_type,required"`
	// ExternalSystem is optional and if provided will override value from deposit asset details
	ExternalSystem  int32  `fig:"external_system"`
	FixedDepositFee uint64 `fig:"fixed_deposit_fee,required"`
	DisableVerify   bool   `fig:"disable_verify"`

	Signer keypair.Full    `fig:"signer,required" mapstructure:"signer"`
	Source keypair.Address `fig:"source,required"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"supervisor":           c.Supervisor,
		"last_processed_block": c.LastProcessedBlock,
		"min_deposit_amount":   c.MinDepositAmount,
		"deposit_asset":        c.DepositAsset,
		"external_system":      c.ExternalSystem,
		"fixed_deposit_fee":    c.FixedDepositFee,
		"disable_verify":       c.DisableVerify,
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
