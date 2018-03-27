package internal

import (
	"bytes"
	"context"
	"encoding/hex"
	"math/big"

	"time"

	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
)

type TxCreator interface {
	CreateTX(tx string, amount int64) (string, error)
}

type ETHHelper struct {
	TxCreator
	log        *logan.Entry
	eth        *ethclient.Client
	address    common.Address
	wallet     *eth.Wallet
	gasPrice   *big.Int
	token      *Token
	marshaller TxMarshaller
}

func NewETHHelper(
	eth *ethclient.Client, address common.Address, wallet *eth.Wallet, gasPrice *big.Int, token *Token,
	log *logan.Entry,
) *ETHHelper {
	var txCreator TxCreator
	if token == nil {
		txCreator = NewETHCreator(gasPrice, eth, address, wallet)
	} else {
		txCreator = NewERC20Creator(eth, token, address, gasPrice)
	}
	return &ETHHelper{
		txCreator,
		log,
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
		if !strings.Contains(err.Error(), "known transaction") {
			return "", errors.Wrap(err, "failed to submit tx")
		}
	}

	// wait while transaction is mined to avoid nonce mismatch issues
	h.ensureMined(context.TODO(), tx.Hash())

	return tx.Hash().Hex(), nil
}

// TODO move this to eth client abstraction
func (h *ETHHelper) ensureMined(ctx context.Context, hash common.Hash) {
	app.RunUntilSuccess(ctx, h.log, "ensure-mined", func(i context.Context) error {
		tx, pending, err := h.eth.TransactionByHash(ctx, hash)
		if err != nil {
			return errors.Wrap(err, "failed to get tx")
		}
		if pending {
			return errors.New("not yet mined")
		}
		if tx == nil {
			return errors.New("transaction not found")
		}
		return nil
	}, 10*time.Second)
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

// TODO
func (h *ETHHelper) GetHash(txHex string) (string, error) {

	panic("Not implemented!")

	//bb, err := hex.DecodeString(txHex)
	//if err != nil {
	//	return "", errors.Wrap(err, "Failed to decode TX bytes from string")
	//}
	//
	//hashBB := sha256.Sum256(bb)
	//return hex.EncodeToString(hashBB[:]), nil
}
