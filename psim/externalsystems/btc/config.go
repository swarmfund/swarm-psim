package btc

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	Source keypair.Address `fig:"source,required"`
	Signer keypair.Full    `fig:"signer,required"`
	// TargetCount of pool entities deployer will try to achieve
	TargetCount uint64 `fig:"target_count,required"`
	// ExternalTypes pool types for which deployer will create entities if needed
	ExternalTypes []string `fig:"external_types,required"`
	HDKey         string   `fig:"hd_key,required"`
	// TODO NetworkType
}

func NewConfig(raw map[string]interface{}) (*Config, error) {
	config := &Config{}
	err := figure.
		Out(config).
		From(raw).
		With(figure.BaseHooks, utils.ETHHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "failed to figure out")
	}
	return config, nil
}
