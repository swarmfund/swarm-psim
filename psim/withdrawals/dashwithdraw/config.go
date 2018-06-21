package dashwithdraw

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	PrivateKey string `fig:"offchain_private_key,required"`

	HotWalletAddress      string `fig:"hot_wallet_address,required"`
	HotWalletScriptPubKey string `fig:"hot_wallet_script_pub_key,required"`
	HotWalletRedeemScript string `fig:"hot_wallet_redeem_script,required"`
	FetchUTXOFrom         uint64 `fig:"fetch_utxo_from,required"` // BlockNumber

	// DustThreshold (in satoshi) value is used to restrict generating new UTXOs lower than this value
	// and to restrict using an existing UTXO with value lower than this Threshold during CoinSelection.
	DustThreshold      int64   `fig:"dust_output_threshold,required"`
	BlocksToBeIncluded uint    `fig:"blocks_to_be_included,required"` // From 2 to 25
	MaxFeePerKB        float64 `fig:"max_fee_per_kb,required"`

	OffchainCurrency   string `fig:"offchain_currency,required"`
	OffchainBlockchain string `fig:"offchain_blockchain,required"`

	MinWithdrawAmount int64 `fig:"min_withdraw_amount,required"`

	SignerKP keypair.Full `fig:"signer,required" mapstructure:"signer"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"hot_wallet_address":       c.HotWalletAddress,
		"hot_wallet_script_pubkey": c.HotWalletScriptPubKey,
		"hot_wallet_redeem_script": c.HotWalletRedeemScript,
		"fetch_utxo_from":          c.FetchUTXOFrom,

		"dust_output_threshold": c.DustThreshold,
		"blocks_to_be_included": c.BlocksToBeIncluded,
		"max_fee_per_kb":        c.MaxFeePerKB,

		"offchain_currency":   c.OffchainCurrency,
		"offchain_blockchain": c.OffchainBlockchain,
		"min_withdraw_amount": c.MinWithdrawAmount,
	}
}

func NewConfig(configData map[string]interface{}) (Config, error) {
	config := Config{}

	err := figure.
		Out(&config).
		From(configData).
		With(figure.BaseHooks, utils.CommonHooks).
		Please()
	if err != nil {
		return Config{}, errors.Wrap(err, "Failed to figure out")
	}

	return config, nil
}
