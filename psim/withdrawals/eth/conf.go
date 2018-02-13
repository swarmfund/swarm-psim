package eth

import (
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/tokend/keypair"
)

type WithdrawConfig struct {
	Source keypair.Address
	Signer keypair.Full
	Asset  string
	// Threshold minimal withdraw amount in gwei
	Threshold int64
	// Key hex encode private key
	Key string
	// GasPrice transaction gas price in gwei
	GasPrice int64
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
	GasPrice int64
	// Token optional token address, configuring it will trigger ERC20 flow
	Token *common.Address
}
