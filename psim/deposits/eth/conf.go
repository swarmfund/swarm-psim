package erc20

import (
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/tokend/keypair"
)

type DepositConfig struct {
	Source        keypair.Address `fig:"source,required"`
	Signer        keypair.Full    `fig:"signer,required"`
	Cursor        uint64          `fig:"cursor,required"`
	Confirmations uint64          `fig:"confirmations"`
	BaseAsset     string          `fig:"base_asset,required"`
	// DepositAsset swarm asset to deposit
	DepositAsset string `fig:"deposit_asset,required"`
	// Token deposit token contract address
	Token common.Address `fig:"token,required"`
}

type VerifyConfig struct {
	Host          string       `fig:"host"`
	Port          int          `fig:"port"`
	Signer        keypair.Full `fig:"signer,required"`
	Confirmations uint64       `fig:"confirmations"`
	// DepositAsset swarm asset to deposit
	DepositAsset string `fig:"deposit_asset,required"`
	// Token deposit token contract address
	Token common.Address `fig:"token,required"`
}
