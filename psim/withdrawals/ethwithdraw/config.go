package ethwithdraw

import (
	"math/big"

	"strings"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	Source keypair.Address `fig:"source,required"`
	// Signer will be used to sign TokenD transactions as well as to access Horizon resources
	Signer keypair.Full `fig:"signer,required"`

	// Key hex encode private key
	PrivateKey string `fig:"private_key,required"`

	MultisigWallet *common.Address `fig:"multisig_wallet_contract_address,required"`
	// Asset specifies TokenD asset code for which service will try to fulfill withdrawal requests
	Asset string `fig:"asset,required"`
	// AssetPrecision is normally 18 (for ETH and most ERC20 tokens), but this field is required
	// in order not to forget to put this value when precision is not 18.
	AssetPrecision uint `fig:"asset_precision,required"`
	// TokenAddress is optional, configuring it will trigger ERC20 flow
	// If not provided - withdraw service will work with ETH
	TokenAddress *common.Address `fig:"token_address"`

	// MinWithdrawAmount minimal withdraw amount with asset precision (in Wei for ETH and most ETC20 tokens)
	MinWithdrawAmount *big.Int `fig:"min_withdraw_amount,required"`
	// GasPrice transaction gas price in Wei
	GasPrice *big.Int `fig:"gas_price,required"`
	GasLimit uint64   `fig:"gas_limit,required"`

	ETHTxsWhiteList []string `fig:"eth_txs_white_list"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"multisig_wallet_contract_address": c.MultisigWallet.String(),
		"asset":              c.Asset,
		"asset_precision":    c.AssetPrecision,
		"token_address":      c.TokenAddress.String(),
		"threshold":          c.MinWithdrawAmount,
		"gas_price":          c.GasPrice.String(),
		"gas_limit":          c.GasLimit,
		"eth_txs_white_list": c.ETHTxsWhiteList,
	}
}

func (c Config) Validate() error {
	if c.Asset != "ETH" && c.TokenAddress == nil {
		return errors.New("Missing TokenAddress - it can only be omitted if Asset is 'ETH'.")
	}

	return nil
}

func (c Config) IsETHTxWhitelisted(ethTXHash string) bool {
	if strings.HasPrefix(ethTXHash, "0x") {
		ethTXHash = ethTXHash[2:]
	}

	for _, h := range c.ETHTxsWhiteList {
		if strings.HasPrefix(h, "0x") {
			h = h[2:]
		}

		if h == ethTXHash {
			return true
		}
	}

	return false
}

func NewConfig(configData map[string]interface{}) (*Config, error) {
	config := &Config{}

	err := figure.
		Out(config).
		With(figure.BaseHooks, utils.ETHHooks).
		From(configData).
		Please()

	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out")

	}

	return config, nil
}
