package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/keypair"
)

type WithdrawConfig struct {
	// Signer will be used to sign TokenD transactions as well as to access Horizon resources
	Signer keypair.Full `fig:"signer,required"`
	// Asset specifies TokenD asset code for which service will try to fulfill withdrawal requests
	Asset string `fig:"asset,required"`
	// Threshold minimal withdraw amount in gwei
	Threshold int64 `fig:"threshold,required"`
	// Key hex encode private key
	Key string `fig:"key,required"`
	// GasPrice transaction gas price in gwei
	GasPrice *big.Int `fig:"gas_price,required"`
	// Token optional token address, configuring it will trigger ERC20 flow
	Token *common.Address `fig:"token"`
	// VerifierServiceName, set the name of which service need to use when verify eth withdraw
	VerifierServiceName string `fig:"verifier_service_name,required"`
}

func NewWithdrawConfig(configData map[string]interface{}) (*WithdrawConfig, error) {
	config := &WithdrawConfig{}

	err := figure.
		Out(config).
		With(figure.BaseHooks, utils.ETHHooks).
		From(configData).
		Please()

	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out")

	}

	return config, nil
}

type VerifyConfig struct {
	Host   string
	Port   int
	Source keypair.Address
	Signer keypair.Full
	Asset  string
	// Threshold minimal withdraw amount in gwei
	Threshold int64
	// Key hex encode private key
	Key string
	// GasPrice transaction gas price in gwei
	GasPrice *big.Int
	// Token optional token address, configuring it will trigger ERC20 flow
	Token *common.Address
}
