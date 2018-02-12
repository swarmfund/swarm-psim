package erc20

import (
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/tokend/keypair"
)

type DepositConfig struct {
	Source        keypair.Address
	Signer        keypair.Full
	Cursor        uint64
	Confirmations uint64
	BaseAsset     string
	// DepositAsset swarm asset to deposit
	DepositAsset string
	// Token deposit token contract address
	Token common.Address
}

type VerifyConfig struct {
	Host          string
	Port          int
	Signer        keypair.Full
	Cursor        uint64
	Confirmations uint64
	// DepositAsset swarm asset to deposit
	DepositAsset string
	// Token deposit token contract address
	Token common.Address
}
