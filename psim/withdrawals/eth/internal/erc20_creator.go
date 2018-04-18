package internal

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

type ERC20Creator struct {
	eth        *ethclient.Client
	token      *Token
	address    common.Address
	gasPrice   *big.Int
	marshaller TxMarshaller
}

func NewERC20Creator(eth *ethclient.Client, token *Token, address common.Address, gasPrice *big.Int) *ERC20Creator {
	return &ERC20Creator{
		eth, token, address, gasPrice, TxMarshaller{},
	}
}

func (h *ERC20Creator) CreateTX(desthex string, amount int64) (string, error) {
	destination := common.HexToAddress(desthex)

	nonce, err := h.eth.PendingNonceAt(context.TODO(), h.address)
	if err != nil {
		return "", errors.Wrap(err, "failed to get nonce")
	}

	// it's here for debug reasons, remove it if you have fresh hot-wallet
	if nonce == 0 {
		return "", errors.Wrap(err, "zero nonce")
	}

	// FIXME withdraw fee is hardcoded in token asset
	input, err := h.token.Transfer(destination, fromGwei(big.NewInt(amount-1250000000)))
	if err != nil {
		return "", errors.Wrap(err, "failed to build tx input")
	}

	tx := types.NewTransaction(
		nonce, h.token.Address(), big.NewInt(0), 200000, fromGwei(h.gasPrice), input)

	// asserting just in case
	if tx.Nonce() != nonce {
		return "", errors.Wrap(err, "geth client is shit")
	}

	return h.marshaller.Marshal(tx)
}
