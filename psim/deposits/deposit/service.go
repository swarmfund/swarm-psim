package deposit

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/issuance"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/keypair"
)

var (
	ErrNoVerifierServices = errors.New("No Deposit Verify services were found.")
)

// AddressProvider must be implemented by WatchAddress storage to pass into Service constructor.
type AddressProvider interface {
	ExternalAccountAt(ctx context.Context, ts time.Time, externalSystem int32, externalData string) (address *string)
	Balance(ctx context.Context, address string, asset string) (balance *string)
}

// Discovery must be implemented by Discovery(Consul) client to pass into Service constructor.
type Discovery interface {
	DiscoverService(service string) ([]discovery.ServiceEntry, error)
}

// OffchainHelper is the interface for specific Offchain(e.g. BTC or ETH)
// deposit and deposit-verify services to implement
// and parametrise the Service.
type OffchainHelper interface {
	// GetAsset must return the name of the Asset being issued in the system during the Deposits processing.
	// Should be configured via config.
	GetAsset() string
	// GetMinDepositAmount must return minimal value for Deposit in Offchain precision.
	GetMinDepositAmount() uint64
	// GetFixedDepositFee must return the value of the fixed Deposit fee in Offchain precision.
	// We substitute Deposit fee for future moving of money to hot wallets.
	//
	// This value is not static and configured because of dynamic fee rates in Offchains.
	GetFixedDepositFee() uint64

	// ConvertToSystem must convert the value with the offchain precision to the system precision.
	// The ONE value from the package amount shows current number of units in one system token.
	ConvertToSystem(offchainAmount uint64) (systemAmount uint64)
	// BuildReference must return a unique identifier of the Deposit, build from Offchain data.
	//
	// Reference is submitted to core and is used to prevent multiple Deposits about the same Offchain TX.
	// Be sure for the same arguments - the reference will always be the same, otherwise duplicate Issuances can appear.
	//
	// You probably won't use all of the provided arguments, but it's no problem
	// for the abstract deposit service to provide all this values into implementations.
	BuildReference(blockNumber uint64, txHash, offchainAddress string, outIndex uint, maxLen int) string
	// GetAddressSynonyms must return all possible representations of the provided Offchain Address
	// which represent the same Address (for example for Ether - lowercased Address is equal to the initial).
	//
	// The returned slice must contain at least 1 element.
	//
	// If returned slice contains only 1 element - it must be the one provided as parameter.
	//
	// The elements in the returned slice *should* not be duplicated.
	GetAddressSynonyms(address string) []string

	// GetLastKnownBlockNumber must return the number of last Block currently existing in the Offchain.
	GetLastKnownBlockNumber() (uint64, error)
	// GetBlock must retrieve the Block with number `number` from the Offchain and parse it into the type Block.
	// It's OK to have Outs in a Tx with Addresses equal to empty string (if output failed to parse or definitely not interesting for us).
	GetBlock(number uint64) (*Block, error)
}

// Service implements app.Service interface.
// Service supervises Offchain blockchain,
// detects arriving deposits,
// verifies Issuance via Verifier
// and sends Issuance(CoinEmissionRequests) to Horizon.
//
// Service uses OffchainHelper to do offchain-specific operations.
type Service struct {
	log *logan.Entry

	source keypair.Address
	signer keypair.Full

	serviceName        string
	lastProcessedBlock uint64
	externalSystem     int32

	// TODO Interface
	horizon         *horizon.Connector
	addressProvider AddressProvider
	builder         *xdrbuild.Builder
	offchainHelper  OffchainHelper
}

// New is constructor for the deposit Service.
//
// Make sure HorizonConnector provided to constructor is with signer.
func New(opts *Opts) *Service {
	return &Service{
		log:                opts.Log.WithField("service", opts.ServiceName),
		source:             opts.Source,
		signer:             opts.Signer,
		serviceName:        opts.ServiceName,
		lastProcessedBlock: opts.LastProcessedBlock,
		horizon:            opts.Horizon,
		addressProvider:    opts.AddressProvider,
		builder:            opts.Builder,
		offchainHelper:     opts.OffchainHelper,
		externalSystem:     opts.ExternalSystem,
	}
}

type Opts struct {
	Log                *logan.Entry
	Source             keypair.Address
	Signer             keypair.Full
	ServiceName        string
	LastProcessedBlock uint64
	Horizon            *horizon.Connector
	ExternalSystem     int32
	AddressProvider    AddressProvider
	Builder            *xdrbuild.Builder
	OffchainHelper     OffchainHelper
}

func (s *Service) Run(ctx context.Context) {
	s.log.Info("Starting.")
	running.WithBackOff(ctx, s.log, s.serviceName, s.processNewBlocks, 5*time.Second, 5*time.Second, time.Hour)
}

func (s *Service) processNewBlocks(ctx context.Context) error {
	lastKnownBlock, err := s.offchainHelper.GetLastKnownBlockNumber()
	if err != nil {
		return errors.Wrap(err, "Failed to GetBlockCount")
	}

	for i := s.lastProcessedBlock + 1; i <= lastKnownBlock; i++ {
		if running.IsCancelled(ctx) {
			return nil
		}

		err := s.processBlock(ctx, i)
		if err != nil {
			return errors.Wrap(err, "Failed to process Block", logan.F{"block_index": i})
		}

		if running.IsCancelled(ctx) {
			// Don't update lastProcessedBlock, because Block processing was probably not finished - ctx was canceled.
			return nil
		}

		s.lastProcessedBlock = i
	}

	return nil
}

func (s *Service) processBlock(ctx context.Context, blockNumber uint64) error {
	s.log.WithField("block_number", blockNumber).Debug("Processing block.")

	block, err := s.offchainHelper.GetBlock(blockNumber)
	if err != nil {
		return errors.Wrap(err, "Failed to get Block from OffchainHelper")
	}

	if block == nil {
		// helper thinks it's 404, we don't care
		return nil
	}

	for _, tx := range block.TXs {
		if running.IsCancelled(ctx) {
			return nil
		}

		err := s.processTX(ctx, blockNumber, block.Timestamp, tx)
		if err != nil {
			return errors.Wrap(err, "Failed to process TX", logan.F{
				"block": block,
				"tx":    tx,
			})
		}
	}

	return nil
}

func (s *Service) processTX(ctx context.Context, blockNumber uint64, blockTime time.Time, tx Tx) error {
	for i, out := range tx.Outs {
		if running.IsCancelled(ctx) {
			return nil
		}

		if out.Address == "" {
			continue
		}

		var accountAddress *string
		addresses := s.offchainHelper.GetAddressSynonyms(out.Address)
		for _, addr := range addresses {
			accountAddress = s.addressProvider.ExternalAccountAt(ctx, blockTime, s.externalSystem, addr)
			if accountAddress != nil {
				// Found
				break
			}

			if running.IsCancelled(ctx) {
				return nil
			}
		}

		if accountAddress == nil {
			// No our Account found for this Offchain Address, skipping.
			continue
		}

		fields := logan.F{
			"out_value":       out.Value,
			"out_index":       i,
			"offchain_addr":   out.Address,
			"account_address": *accountAddress,
		}

		if out.Value < s.offchainHelper.GetMinDepositAmount() {
			s.log.WithFields(fields.Merge(logan.F{
				"block_number": blockNumber,
				"block_time":   blockTime,
				"tx_hash":      tx.Hash,
			})).WithField("min_deposit_amount_from_config", s.offchainHelper.GetMinDepositAmount()).
				Warn("Received deposit with too small amount.")
			continue
		}

		err := s.processDeposit(ctx, blockNumber, blockTime, tx.Hash, uint(i), out, *accountAddress)
		// TODO Retry processing Deposit - don't go to following Deposits.
		if err != nil {
			return errors.Wrap(err, "Failed to process Deposit", fields)
		}
	}

	return nil
}

func (s *Service) processDeposit(ctx context.Context, blockNumber uint64, blockTime time.Time, txHash string, outIndex uint, out Out, accountAddress string) error {
	fields := logan.F{
		"block_number":    blockNumber,
		"tx_hash":         txHash,
		"out_index":       outIndex,
		"offchain_addr":   out.Address,
		"offchain_amount": out.Value,
		"account_address": accountAddress,
	}
	s.log.WithFields(fields).Debug("Processing deposit.")

	balanceID := s.addressProvider.Balance(ctx, accountAddress, s.offchainHelper.GetAsset())
	if balanceID == nil {
		// user does not have target balance
		// unfortunate, but we don't care
		s.log.WithFields(fields).Warn("no deposit asset balance found")
		return nil
	}

	valueWithoutDepositFee := out.Value - s.offchainHelper.GetFixedDepositFee()
	emissionAmount := s.offchainHelper.ConvertToSystem(valueWithoutDepositFee)

	reference := s.offchainHelper.BuildReference(blockNumber, txHash, out.Address, outIndex, 64)
	// TODO check maxLen

	issuanceOpt := issuance.RequestOpt{
		Reference: reference,
		Receiver:  *balanceID,
		Asset:     s.offchainHelper.GetAsset(),
		Amount:    emissionAmount,
		Details: ExternalDetails{
			BlockNumber: blockNumber,
			TXHash:      txHash,
			OutIndex:    outIndex,
			Price:       amount.One,
		}.Encode(),
		// TODO stop setting this Task in future, it should be set as default Task for IssuanceRequest in Core in future.
		AllTasks: issuance.TaskVerifyDeposit,
	}

	fields = fields.Merge(logan.F{
		"balance_id":              balanceID,
		"converted_system_amount": emissionAmount,
		"reference":               reference,
	})

	err := s.processIssuance(ctx, blockNumber, out.Address, accountAddress, issuanceOpt)
	if err != nil {
		return errors.Wrap(err, "Failed to process Issuance", fields)
	}

	return nil
}

func (s *Service) processIssuance(ctx context.Context, blockNumber uint64, offchainAddress, accountAddress string, issuanceOpt issuance.RequestOpt) error {
	tx := issuance.CraftIssuanceTX(issuanceOpt, s.builder, s.source, s.signer)

	envelope, err := tx.Marshal()
	if err != nil {
		return errors.Wrap(err, "Failed to marshal TX into envelope")
	}

	logger := s.log.WithFields(logan.F{
		"block_number":     blockNumber,
		"offchain_address": offchainAddress,
		"account_address":  accountAddress,
		"issuance":         issuanceOpt,
	})

	ok, err := issuance.SubmitEnvelope(ctx, envelope, s.horizon.Submitter())
	if err != nil {
		logger.WithError(err).Error("Failed to submit CoinEmissionRequest TX to Horizon.")
		return nil
	}

	if ok {
		logger.Info("CoinEmissionRequest was sent successfully.")
	} else {
		logger.Debug("Reference duplication - already processed Deposit, skipping.")
	}
	return nil
}
