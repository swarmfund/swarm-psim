package ethcontracts

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	Source keypair.Address `fig:"source,required"`
	Signer keypair.Full    `fig:"signer,required"`
	// TargetCount of pool entities deployer will try to achieve
	TargetCount uint64 `fig:"target_count,required"`
	// ExternalTypes pool types for which deployer will create entities if needed
	ExternalTypes []string       `fig:"external_types,required"`
	ETHPrivateKey string         `fig:"eth_private_key,required"`
	ContractOwner common.Address `fig:"contract_owner,required"`
	// GasPrice in gwei
	GasPrice *big.Int `fig:"gas_price,required"`
	// GasLimit in gwei
	GasLimit *big.Int `fig:"gas_limit"`
}

func NewConfig(raw map[string]interface{}) (*Config, error) {
	config := &Config{
		GasLimit: big.NewInt(420000),
	}
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
