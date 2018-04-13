package btcwithdraw

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	PrivateKey string `fig:"btc_private_key,required"`

	HotWalletAddress      string `fig:"hot_wallet_address,required"`
	HotWalletScriptPubKey string `fig:"hot_wallet_script_pub_key,required"`
	HotWalletRedeemScript string `fig:"hot_wallet_redeem_script,required"`

	MinWithdrawAmount int64 `fig:"min_withdraw_amount,required"`

	SignerKP keypair.Full `fig:"signer,required" mapstructure:"signer"`
}

func NewConfig(configData map[string]interface{}) (*Config, error) {
	config := &Config{}

	err := figure.
		Out(config).
		From(configData).
		With(figure.BaseHooks, utils.CommonHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out", logan.F{
			"service": conf.ServiceBTCWithdraw,
		})
	}

	return config, nil
}
