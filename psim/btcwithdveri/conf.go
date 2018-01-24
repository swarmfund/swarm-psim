package btcwithdveri

import "gitlab.com/tokend/keypair"

type Config struct {
	Host string `fig:"host"`
	Port int    `fig:"port"`
	// TODO Pprof field?

	PrivateKey string `fig:"btc_private_key"`

	HotWalletAddress string `fig:"hot_wallet_address"`
	HotWalletScriptPubKey string `fig:"hot_wallet_script_pub_key"`
	HotWalletRedeemScript string `fig:"hot_wallet_redeem_script"`

	MinWithdrawAmount float64 `fig:"min_withdraw_amount"`

	SourceKP keypair.Address `fig:"source"`
	SignerKP keypair.Full    `fig:"signer" mapstructure:"signer"`
}
