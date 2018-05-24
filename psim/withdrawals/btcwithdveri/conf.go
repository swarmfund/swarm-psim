package btcwithdveri

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	Host string `fig:"host"`
	Port int    `fig:"port"`
	// TODO Pprof field?

	PrivateKey string `fig:"btc_private_key,required"`

	HotWalletAddress      string `fig:"hot_wallet_address,required"`
	HotWalletScriptPubKey string `fig:"hot_wallet_script_pub_key,required"`
	HotWalletRedeemScript string `fig:"hot_wallet_redeem_script,required"`

	OffchainCurrency   string `fig:"offchain_currency,required"`
	OffchainBlockchain string `fig:"offchain_blockchain,required"`

	MinWithdrawAmount int64 `fig:"min_withdraw_amount,required"`

	SourceKP keypair.Address `fig:"source,required"`
	SignerKP keypair.Full    `fig:"signer,required"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"host":        c.Host,
		"port":        c.Port,

		"hot_wallet_address":        c.HotWalletAddress,
		"hot_wallet_script_pub_key": c.HotWalletScriptPubKey,
		"hot_wallet_redeem_script":  c.HotWalletRedeemScript,

		"offchain_currency":         c.OffchainCurrency,
		"offchain_blockchain":       c.OffchainBlockchain,

		"min_withdraw_amount":       c.MinWithdrawAmount,
	}
}

func NewConfig(configData map[string]interface{}) (*Config, error) {
	config := &Config{}

	err := figure.
		Out(config).
		From(configData).
		With(figure.BaseHooks, utils.CommonHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out")
	}

	return config, nil
}
