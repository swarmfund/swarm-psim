package eth_deposit

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// TODO consider change amount type

type TxProvider interface {
	Send(amount int, destination string) (succes bool, err error)
}

type EthTxProvider struct {
	eclient *ethclient.Client
	ctx     context.Context
	kp      eth.Keypair
	log     *logan.Entry
}

func NewEthTxProvider(eclient *ethclient.Client, ctx context.Context, kp eth.Keypair, log *logan.Entry) TxProvider {
	return &EthTxProvider{
		eclient,
		ctx,
		kp,
		log,
	}
}

func (e *EthTxProvider) Send(amount int, externalAddress string) (bool, error) {
	// transfer some ETH
	nonce, err := e.eclient.NonceAt(e.ctx, e.kp.Address(), nil)
	if err != nil {
		return false, errors.Wrap(err, "failed to get address nonce")
	}

	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(externalAddress),
		eth.FromGwei(big.NewInt(2000)),
		22000,
		eth.FromGwei(big.NewInt(5)),
		nil,
	)

	tx, err = e.kp.SignTX(tx)
	if err != nil {
		return false, errors.Wrap(err, "failed to sign tx")
	}

	if err = e.eclient.SendTransaction(e.ctx, tx); err != nil {
		return false, errors.Wrap(err, "failed to send transfer tx")
	}

	eth.EnsureHashMined(e.ctx, e.log, e.eclient, tx.Hash())
	return true, nil
}