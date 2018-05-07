package internal

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
	"gitlab.com/tokend/go/amount"
)

type Helper struct {
	*Converter
	*ConfigHelper
	*ETHHelper
}

func NewHelper(
	asset string, withdrawThreshold int64, eth *ethclient.Client, address common.Address, wallet *eth.Wallet,
	gasPrice *big.Int, token *Token, log *logan.Entry,
) *Helper {
	return &Helper{
		NewConverter(),
		NewConfigHelper(asset, withdrawThreshold),
		NewETHHelper(eth, address, wallet, gasPrice, token, log),
	}
}

type Converter struct {
}

func NewConverter() *Converter {
	return &Converter{}
}

func (h *Converter) ConvertAmount(dest int64) int64 {
	// expected offchain to be in gwei precision (10^9)
	var gwei int64 = 1000000000
	result, overflow := amount.BigDivide(dest, gwei, amount.One, amount.ROUND_DOWN)
	if overflow {
		panic("overflow")
	}
	return result
}

func fromGwei(amount *big.Int) *big.Int {
	return new(big.Int).Mul(amount, new(big.Int).SetInt64(1000000000))
}

func toGwei(amount *big.Int) *big.Int {
	return new(big.Int).Div(amount, new(big.Int).SetInt64(1000000000))
}

type ConfigHelper struct {
	asset     string
	threshold int64
}

func NewConfigHelper(asset string, threshold int64) *ConfigHelper {
	return &ConfigHelper{
		asset:     asset,
		threshold: threshold,
	}
}

func (h *ConfigHelper) GetAsset() string {
	return h.asset
}

func (h *ConfigHelper) GetMinWithdrawAmount() int64 {
	return h.threshold
}
