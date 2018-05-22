package ethwithdraw

import (
	"context"

	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/horizon-connector"
)

type WithdrawRequestsStreamer interface {
	StreamWithdrawalRequestsOfAsset(ctx context.Context, destAssetCode string, reverseOrder, endlessly bool) <-chan horizon.ReviewableRequestEvent
}

type TXSubmitter interface {
	SubmitE(txEnvelope string) (horizon.SubmitResponseDetails, error)
}

type ETHClient interface {
	bind.ContractBackend
	//SendTransaction(ctx context.Context, tx *types.Transaction) error
	//PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)

	TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
}

type ETHWallet interface {
	SignTX(address common.Address, tx *types.Transaction) (*types.Transaction, error)
}

type Service struct {
	log        *logan.Entry
	config     Config
	ethAddress common.Address

	withdrawRequestsStreamer WithdrawRequestsStreamer
	xdrbuilder               *xdrbuild.Builder
	txSubmitter              TXSubmitter

	ethClient        ETHClient
	ethWallet        ETHWallet
	multisigContract *eth.MultisigWalletTransactor

	newETHSequence uint64
}

func NewService(
	log *logan.Entry,
	config Config,
	ethAddress common.Address,
	streamer WithdrawRequestsStreamer,
	xdrbuilder *xdrbuild.Builder,
	txSubmitter TXSubmitter,
	ethClient ETHClient,
	ethWallet ETHWallet) (*Service, error) {

	multisigContract, err := eth.NewMultisigWalletTransactor(*config.MultisigWallet, eth.NewContractTransactor(ethClient))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create MultisigWallet Contract")
	}

	return &Service{
		log:        log.WithField("service_name", conf.ServiceETHWithdraw),
		config:     config,
		ethAddress: ethAddress,

		withdrawRequestsStreamer: streamer,
		xdrbuilder:               xdrbuilder,
		txSubmitter:              txSubmitter,

		ethClient:        ethClient,
		ethWallet:        ethWallet,
		multisigContract: multisigContract,
	}, nil
}

func (s *Service) Run(ctx context.Context) {
	s.log.WithField("", s.config).Info("Started.")

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		s.submitETHTransactionsInfinitely(ctx)
		wg.Done()
	}()

	var err error
	// blocking call here intentionally - we don't start approving requests until get proper ETH sequence
	s.newETHSequence, err = s.detectNewETHSequence(ctx)
	if err != nil {
		s.log.WithError(err).Error("Failed to detect new ETH TX Sequence - critical error, stopping.")
		return
	}
	if running.IsCancelled(ctx) {
		wg.Wait()
		s.log.Info("Service stopped smoothly during detecting of ETH sequence(nonce), " +
			"without starting approving/rejecting requests.")
		return
	}

	wg.Add(1)
	go func() {
		s.processTSWRequestsInfinitely(ctx)
		wg.Done()
	}()

	wg.Wait()
	s.log.Info("Service stopped smoothly.")
}
