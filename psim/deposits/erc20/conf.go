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
	// DepositAsset TokenD asset to deposit
	DepositAsset string `fig:"deposit_asset,required"`
	// Token deposit token contract address
	Token common.Address `fig:"token"`
	// ExternalSystem type used for matching deposit addresses,
	// if set will override one in deposit asset details
	ExternalSystem int32 `fig:"external_system"`
	DisableVerify  bool  `fig:"disable_verify"`
}

type VerifyConfig struct {
	Host          string       `fig:"host"`
	Port          int          `fig:"port"`
	Signer        keypair.Full `fig:"signer,required"`
	Confirmations uint64       `fig:"confirmations"`
	// DepositAsset swarm asset to deposit
	DepositAsset   string `fig:"deposit_asset,required"`
	ExternalSystem int32  `fig:"external_system"`
	// Token deposit token contract address
	Token common.Address `fig:"token,required"`
}
