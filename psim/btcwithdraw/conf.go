package btcwithdraw

import "gitlab.com/swarmfund/go/keypair"

type Config struct {
	PrivateKey string `fig:"private_key"`
	// TODO Remove after implementing of verifier
	PrivateKey2      string `fig:"private_key_2"`
	HotWalletAddress string `fig:"hot_wallet_address"`

	SourceKP keypair.KP `fig:"exchange"`
	SignerKP keypair.KP `fig:"signer" mapstructure:"signer"`
}
