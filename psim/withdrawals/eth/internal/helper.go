package internal

import (
	"context"
	"math/big"

	"bytes"

	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
)

type Helper struct {
	*Converter
	*ConfigHelper
	*ETHHelper
}

func NewHelper(
	asset string, withdrawThreshold int64, eth *ethclient.Client, address common.Address, wallet *eth.Wallet,
	gasPrice int64, token *Token,
) *Helper {
	return &Helper{
		NewConverter(),
		NewConfigHelper(asset, withdrawThreshold),
		NewETHHelper(eth, address, wallet, gasPrice, token),
	}
}

type ETHHelper struct {
	eth      *ethclient.Client
	address  common.Address
	wallet   *eth.Wallet
	gasPrice int64
	token    *Token
}

func NewETHHelper(
	eth *ethclient.Client, address common.Address, wallet *eth.Wallet, gasPrice int64, token *Token,
) *ETHHelper {
	return &ETHHelper{
		eth,
		address,
		wallet,
		gasPrice,
		token,
	}
}

func (h *ETHHelper) ValidateAddress(addr string) error {
	if !common.IsHexAddress(addr) {
		return errors.New("not a valid eth address")
	}
	return nil
}

func (h *ETHHelper) CreateTX(desthex string, amount int64) (string, error) {
	destination := common.HexToAddress(desthex)

	nonce, err := h.eth.PendingNonceAt(context.TODO(), h.address)
	if err != nil {
		return "", errors.Wrap(err, "failed to get nonce")
	}

	input, err := h.token.Transfer(destination, fromGwei(big.NewInt(amount)))
	if err != nil {
		return "", errors.Wrap(err, "failed to build tx input")
	}

	tx := types.NewTransaction(
		nonce, h.token.Address(), big.NewInt(0), big.NewInt(200000), fromGwei(big.NewInt(h.gasPrice)), input)

	var buf bytes.Buffer
	if err := tx.EncodeRLP(&buf); err != nil {
		return "", errors.Wrap(err, "failed to encode tx")
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

func (h *ETHHelper) SendTX(txhex string) (hash string, err error) {
	rlpbytes, err := hex.DecodeString(txhex)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode tx hex")
	}
	var tx types.Transaction
	err = tx.DecodeRLP(rlp.NewStream(bytes.NewReader(rlpbytes), 0))
	if err != nil {
		return "", errors.Wrap(err, "failed to decode tx rlp")
	}

	if err = h.eth.SendTransaction(context.TODO(), &tx); err != nil {
		return "", errors.Wrap(err, "failed to submit tx")
	}

	// TODO wait while mined

	return tx.Hash().Hex(), nil
}

func (h *ETHHelper) SignTX(txhex string) (string, error) {
	rlpbytes, err := hex.DecodeString(txhex)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode tx hex")
	}
	tx := &types.Transaction{}
	err = tx.DecodeRLP(rlp.NewStream(bytes.NewReader(rlpbytes), 0))
	if err != nil {
		return "", errors.Wrap(err, "failed to decode tx rlp")
	}
	tx, err = h.wallet.SignTX(h.address, tx)
	if err != nil {
		return "", errors.Wrap(err, "failed to sign tx")
	}
	var buf bytes.Buffer
	if err := tx.EncodeRLP(&buf); err != nil {
		return "", errors.Wrap(err, "failed to encode tx")
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

func (h *ETHHelper) ValidateTX(tx string, withdrawAddress string, withdrawAmount int64) (string, error) {
	// FIXME currently we are just mimicking real two-step flow
	return "", nil
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
		asset, threshold,
	}
}

func (h *ConfigHelper) GetAsset() string {
	return h.asset
}

func (h *ConfigHelper) GetMinWithdrawAmount() int64 {
	return h.threshold
}
