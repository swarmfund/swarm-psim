package btcwithdveri

import "gitlab.com/swarmfund/go/keypair"

type Config struct {
	Host string `fig:"host"`
	Port int    `fig:"port"`
	// TODO Pprof field?

	PrivateKey string `fig:"private_key"`

	HotWalletAddress string `fig:"hot_wallet_address"`
	HotWalletScriptPubKey string `fig:"hot_wallet_script_pub_key"`
	HotWalletRedeemScript string `fig:"hot_wallet_redeem_script"`

	MinWithdrawAmount float64 `fig:"min_withdraw_amount"`

	SourceKP keypair.KP `fig:"source"`
	SignerKP keypair.KP `fig:"signer" mapstructure:"signer"`
}
