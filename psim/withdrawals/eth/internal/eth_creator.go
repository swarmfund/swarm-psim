package internal

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
)

type ETHCreator struct {
	gasPrice   *big.Int
	eth        *ethclient.Client
	address    common.Address
	wallet     *eth.Wallet
	marshaller TxMarshaller
}

func NewETHCreator(gasPrice *big.Int, eth *ethclient.Client, address common.Address, wallet *eth.Wallet) *ETHCreator {
	return &ETHCreator{
		gasPrice,
		eth,
		address,
		wallet,
		TxMarshaller{},
	}
}

func (h *ETHCreator) CreateTX(desthex string, amount int64) (string, error) {
	txGas := big.NewInt(21000)
	txFee := new(big.Int).Mul(txGas, h.gasPrice)
	withdrawAmount := fromGwei(big.NewInt(amount))
	destination := common.HexToAddress(desthex)

	nonce, err := h.eth.PendingNonceAt(context.TODO(), h.address)
	if err != nil {
		return "", errors.Wrap(err, "failed to get nonce")
	}

	value := new(big.Int).Sub(withdrawAmount, txFee)

	tx, err := h.wallet.SignTX(
		h.address,
		types.NewTransaction(
			nonce,
			destination,
			value,
			txGas,
			h.gasPrice,
			nil,
		),
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to sign tx")
	}

	return h.marshaller.Marshal(tx)
}
