package btcwithdraw

import "gitlab.com/swarmfund/go/keypair"

type Config struct {
	PrivateKey string `fig:"private_key"`
	// TODO Remove after implementing of verifier
	PrivateKey2      string `fig:"private_key_2"`

	HotWalletAddress string `fig:"hot_wallet_address"`
	HotWalletScriptPubKey string `fig:"hot_wallet_script_pub_key"`
	HotWalletRedeemScript string `fig:"hot_wallet_redeem_script"`

	SourceKP keypair.KP `fig:"source"`
	SignerKP keypair.KP `fig:"signer" mapstructure:"signer"`
}
