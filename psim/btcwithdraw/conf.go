package btcwithdraw

import "gitlab.com/tokend/keypair"

type Config struct {
	PrivateKey string `fig:"private_key"`

	HotWalletAddress      string `fig:"hot_wallet_address"`
	HotWalletScriptPubKey string `fig:"hot_wallet_script_pub_key"`
	HotWalletRedeemScript string `fig:"hot_wallet_redeem_script"`

	MinWithdrawAmount float64 `fig:"min_withdraw_amount"`

	SignerKP keypair.Full `fig:"signer" mapstructure:"signer"`
}
