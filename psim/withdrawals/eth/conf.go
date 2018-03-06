package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/tokend/keypair"
)

type WithdrawConfig struct {
	// Signer will be used to sign TokenD transactions as well as to access Horizon resources
	Signer keypair.Full
	// Asset specifies TokenD asset code for which service will try to fulfill withdrawal requests
	Asset string
	// Threshold minimal withdraw amount in gwei
	Threshold int64
	// Key hex encode private key
	Key string
	// GasPrice transaction gas price in gwei
	GasPrice *big.Int
	// Token optional token address, configuring it will trigger ERC20 flow
	Token *common.Address
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
