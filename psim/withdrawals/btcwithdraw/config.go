package btcwithdraw

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/externalsystems/derive"
	"gitlab.com/swarmfund/psim/psim/supervisor"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	PrivateKey            string `fig:"btc_private_key,required"`
	HotWalletAddress      string `fig:"hot_wallet_address,required"`
	HotWalletScriptPubKey string `fig:"hot_wallet_script_pub_key,required"`
	HotWalletRedeemScript string `fig:"hot_wallet_redeem_script,required"`
	// NetworkType will be used to configure chain params
	NetworkType       derive.NetworkType `fig:"network_type,required"`
	MinWithdrawAmount int64              `fig:"min_withdraw_amount,required"`
	// DepositAsset TokenD asset code to issue
	DepositAsset string `fig:"deposit_asset,required"`
	// SignerKP used to access Horizon resources and sign transactions
	SignerKP keypair.Full `fig:"signer,required" mapstructure:"signer"`
	// VerifierServiceName used to discovery verifier service
	VerifierServiceName string `fig:"verifier_service_name,required"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"hot_wallet_address":        c.HotWalletAddress,
		"hot_wallet_script_pub_key": c.HotWalletScriptPubKey,
		"hot_wallet_redeem_script":  c.HotWalletRedeemScript,
		"network_type":              c.NetworkType,
		"min_withdraw_amount":       c.MinWithdrawAmount,
	}
}

func NewConfig(configData map[string]interface{}) (*Config, error) {
	config := &Config{}

	err := figure.
		Out(config).
		From(configData).
		With(figure.BaseHooks, utils.CommonHooks, supervisor.DLFigureHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out")
	}

	return config, nil
}
