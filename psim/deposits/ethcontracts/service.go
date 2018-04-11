package ethcontracts

import (
	"context"
	"math/big"

	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
)

type ETHClient interface {
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
}

type ETHWallet interface {
	SignTX(address common.Address, tx *types.Transaction) (*types.Transaction, error)
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
	//contractAddr, err := s.deployContract(ctx)
	//if err != nil {
	//	s.log.WithError(err).Error("Failed to deploy Contract.")
	//}
	//if running.IsCancelled(ctx) {
	//	return
	//}
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

	// FIXME GasLimit to config
	tx := types.NewContractCreation(nonce, big.NewInt(0), big.NewInt(407266), s.config.GasPrice, contractBytes)
	tx, err = s.ethWallet.SignTX(s.ethAddress, tx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to sign ContractCreation TX", fields)
	}
	fields["tx_hash"] = tx.Hash()

	err = s.ethClient.SendTransaction(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send Transaction", fields)
	}
	if running.IsCancelled(ctx) {
		return nil, nil
	}

	// context.Background() is used intentionally - we don't stop on ctx cancel until submit the newly created Contract to Core.
	s.ensureMined(context.Background(), tx.Hash())

	// context.Background() is used intentionally - we don't stop on ctx cancel until submit the newly created Contract to Core.
	receipt, err := s.ethClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get TransactionReceipt of the ContactCreation TX just been mined", fields)
	}

	return &receipt.ContractAddress, nil
}

// EnsureMined will only return if TX is mined or ctx is cancelled.
func (s *Service) ensureMined(ctx context.Context, hash common.Hash) {
	running.UntilSuccess(ctx, s.log.WithField("tx_hash", hash), "tx_mined_ensurer", func(ctx context.Context) (bool, error) {
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
