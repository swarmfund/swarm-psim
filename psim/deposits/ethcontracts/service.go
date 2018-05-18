package ethcontracts

import (
	"context"
	"math/big"

	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

type ETHClient interface {
	bind.ContractBackend
	//SendTransaction(ctx context.Context, tx *types.Transaction) error

	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
}

type ETHWallet interface {
	SignTX(address common.Address, tx *types.Transaction) (*types.Transaction, error)
	Addresses(ctx context.Context) (result []common.Address)
}

type Service struct {
	log    *logan.Entry
	config Config

	ethAddress common.Address
	ethClient  ETHClient
	ethWallet  ETHWallet
}

func NewService(
	log *logan.Entry,
	config Config,
	ethAddress common.Address,
	ethClient ETHClient,
	ethWallet ETHWallet) *Service {

	return &Service{
		log:    log,
		config: config,

		ethAddress: ethAddress,
		ethClient:  ethClient,
		ethWallet:  ethWallet,
	}
}

func (s *Service) Run(ctx context.Context) {
	contractAddr, err := s.deployContract(ctx)
	if err != nil {
		s.log.WithError(err).Error("Failed to deploy Contract.")
		return
	}
	if running.IsCancelled(ctx) {
		return
	}

	s.log.WithField("contract_address", contractAddr.String()).Info("Contract was deployed successfully.")
}

// DeployContract deploys the contract from the contractBytes package-level var to ETH network
// and waits for the Contract to be deployed (TX which deploys the Contract to be mined).
// If returned error is nil - always returns non-nil Address.
//
// nil, nil will be returned if ctx is cancelled.
func (s *Service) deployContract(ctx context.Context) (*common.Address, error) {
	nonce, err := s.ethClient.NonceAt(ctx, s.ethAddress, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Nonce for the Address")
	}
	if running.IsCancelled(ctx) {
		return nil, nil
	}
	fields := logan.F{
		"current_nonce": nonce,
	}

	contractAddr, tx, _, err := DeployContract(&bind.TransactOpts{
		// TODO Move this potentially paniccing part into constructor.
		From:  s.ethWallet.Addresses(ctx)[0],
		Nonce: big.NewInt(int64(nonce)),
		Signer: func(signer types.Signer, addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
			return s.ethWallet.SignTX(addr, tx)
		},
		Value:    big.NewInt(0),
		GasPrice: s.config.GasPrice,
		// FIXME GasLimit to config
		GasLimit: 500000,
		Context:  ctx,
	}, s.ethClient, s.config.ContractOwner)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build Contract with Owner from Config", fields)
	}
	if running.IsCancelled(ctx) {
		return nil, nil
	}
	fields["contract_address"] = contractAddr.String()
	fields["tx_hash"] = tx.Hash().String()

	// context.Background() is used intentionally - we don't stop on ctx cancel until submit the newly created Contract to ETH network.
	s.ensureMined(context.Background(), tx.Hash())

	// context.Background() is used intentionally - we don't stop on ctx cancel until submit the newly created Contract to ETH network.
	receipt, err := s.ethClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get TransactionReceipt of the ContactCreation TX just been mined", fields)
	}

	return &receipt.ContractAddress, nil
}

// EnsureMined will only return if TX is mined or ctx is cancelled.
func (s *Service) ensureMined(ctx context.Context, hash common.Hash) {
	running.UntilSuccess(ctx, s.log.WithField("tx_hash", hash.String()), "tx_mined_ensurer", func(ctx context.Context) (bool, error) {
		tx, isPending, err := s.ethClient.TransactionByHash(ctx, hash)
		if err != nil {
			return false, errors.Wrap(err, "Failed to get TX by hash from ETHClient")
		}
		if tx == nil {
			return false, errors.New("TX was not found by hash over ETHClient")
		}
		if isPending {
			return false, nil
		}

		return true, nil
	}, 5*time.Second, time.Minute)
}
