package erc20

import "gitlab.com/tokend/keypair"

type DepositConfig struct {
	Source        keypair.Address
	Signer        keypair.Full
	Cursor        uint64
	Confirmations uint64
	BaseAsset     string
	DepositAsset  string
}

type VerifyConfig struct {
	Host          string
	Port          int
	Signer        keypair.Full
	Cursor        uint64
	Confirmations uint64
	DepositAsset  string
}
