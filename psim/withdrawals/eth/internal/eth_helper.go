package internal

import (
	"bytes"
	"context"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
)

type TxCreator interface {
	CreateTX(tx string, amount int64) (string, error)
}

type ETHHelper struct {
	TxCreator
	eth        *ethclient.Client
	address    common.Address
	wallet     *eth.Wallet
	gasPrice   *big.Int
	token      *Token
	marshaller TxMarshaller
}

func NewETHHelper(
	eth *ethclient.Client, address common.Address, wallet *eth.Wallet, gasPrice *big.Int, token *Token,
) *ETHHelper {
	var txCreator TxCreator
	if token == nil {
		txCreator = NewETHCreator(gasPrice, eth, address, wallet)
	} else {
		txCreator = NewERC20Creator(eth, token, address, gasPrice)
	}
	return &ETHHelper{
		txCreator,
		eth,
		address,
		wallet,
		gasPrice,
		token,
		TxMarshaller{},
	}
}

func (h *ETHHelper) ValidateAddress(addr string) error {
	if !common.IsHexAddress(addr) {
		return errors.New("not a valid eth address")
	}
	return nil
}

func (h *ETHHelper) SendTX(txhex string) (hash string, err error) {
	tx, err := h.marshaller.Unmarshal(txhex)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal tx")
	}

	if err = h.eth.SendTransaction(context.TODO(), tx); err != nil {
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
