package eth_deposit

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
	"gitlab.com/tokend/go/xdrbuild"
)

type TxProvider interface {
	Send(ctx context.Context, amount int64, destination string) (succes bool, err error)
	GetWithdrawRequestDetails() xdrbuild.WithdrawRequestDetails
	GetCurrentBalance(ctx context.Context) (*big.Int, error)
}

type EthTxProvider struct {
	eclient *ethclient.Client
	kp      eth.Keypair
	log     *logan.Entry
}

func NewEthTxProvider(eclient *ethclient.Client, kp eth.Keypair, log *logan.Entry) TxProvider {
	return &EthTxProvider{
		eclient,
		kp,
		log,
	}
}

func (e *EthTxProvider) Send(ctx context.Context, amount int64, externalAddress string) (bool, error) {
	// transfer some ETH
	nonce, err := e.eclient.NonceAt(ctx, e.kp.Address(), nil)
	if err != nil {
		return false, errors.Wrap(err, "failed to get address nonce")
	}

	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(externalAddress),
		eth.FromGwei(big.NewInt(amount)),
		22000,
		eth.FromGwei(big.NewInt(5)),
		nil,
	)

	tx, err = e.kp.SignTX(tx)
	if err != nil {
		return false, errors.Wrap(err, "failed to sign tx")
	}

	if err = e.eclient.SendTransaction(ctx, tx); err != nil {
		return false, errors.Wrap(err, "failed to send transfer tx")
	}

	eth.EnsureHashMined(ctx, e.log, e.eclient, tx.Hash())
	return true, nil
}

func (e *EthTxProvider) GetWithdrawRequestDetails() xdrbuild.WithdrawRequestDetails {
	return &xdrbuild.ETHWithdrawRequestDetails{
		Address: e.kp.Address().Hex(),
	}
}

func (e *EthTxProvider) GetCurrentBalance(ctx context.Context) (*big.Int, error) {
	return e.eclient.BalanceAt(ctx, e.kp.Address(), nil)
}
