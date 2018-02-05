package deposit

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/internal/resources"
	"gitlab.com/tokend/keypair"
)

var ErrNoBalanceID = errors.New("No BalanceID for the Account.")

// AddressProvider must be implemented by WatchAddress storage to pass into Service constructor.
type AddressProvider interface {
	AddressAt(ctx context.Context, t time.Time, btcAddress string) (tokendAddress *string)
}

// OffchainHelper is the interface for specific Offchain(BTC or ETH)
// deposit and deposit-verify services to implement
// and parametrise the Service.
type OffchainHelper interface {
	// GetLastKnownBlockNumber must return the number of last Block currently existing in the Offchain.
	GetLastKnownBlockNumber() (uint64, error)
	// GetBlock must retrieve the Block with number `number` from the Offchain and parse it into the type Block.
	// It's OK to have Outs in a Tx with Addresses equal to empty string (if output failed to parse or definitely not interesting for us).
	GetBlock(number uint64) (*Block, error)
	// GetMinDepositAmount must return minimal value for Deposit in Offchain precision.
	GetMinDepositAmount() uint64
	// GetFixedDepositFee must return the value of the fixed Deposit fee in Offchain precision.
	// We substitute Deposit fee for future moving of money to hot wallets.
	// This value is not static and configured due to different fee rates in Offchains.
	GetFixedDepositFee() uint64
	// ConvertToSystem must convert the value with the offchain precision to the system precision.
	// The ONE value from the package amount shows current number of units in one system token.
	ConvertToSystem(offchainAmount uint64) (systemAmount uint64)
	// GetAsset must return the name of the Asset being issued in the system during the Deposits processing.
	// Should be configured via config.
	GetAsset() string
	// BuildReference must return a unique identifier of the Deposit, build from Offchain data.
	// Reference is submitted to core and is used to prevent multiple Deposits about the same Offchain TX.
	// You probably won't use all of the provided arguments, but it's no problem
	// for the abstract deposit service to provide all this values into implementations.
	BuildReference(blockHash, txHash, offchainAddress string, outIndex uint, amount uint64, maxLen int) string
}

// Service implements app.Service interface, it supervises Offchain blockchain
// and send CoinEmissionRequests to Horizon if arrived Deposit detected.
type Service struct {
	log *logan.Entry

	source keypair.Address
	signer keypair.Full
	// TODO Remove after moving logic of the second signature to the verifier.
	additionalSigner keypair.Full

	serviceName        string
	lastProcessedBlock uint64
	lastBlocksNotWatch uint64

	// TODO Interface
	horizon         *horizon.Connector
	addressProvider AddressProvider
	builder         *xdrbuild.Builder
	offchainHelper  OffchainHelper
}

// New is constructor for the deposit Service.
func New(
	log *logan.Entry,
	source keypair.Address,
	signer keypair.Full,
	additionalSigner keypair.Full,
	serviceName string,
	lastProcessedBlock,
	lastBlocksNotWatch uint64,
	horizon *horizon.Connector,
	addressProvider AddressProvider,
	builder *xdrbuild.Builder,
	offchainHelper OffchainHelper) *Service {

	result := &Service{
		log:                log,
		source:             source,
		signer:             signer,
		additionalSigner:   additionalSigner,
		serviceName:        serviceName,
		lastProcessedBlock: lastProcessedBlock,
		lastBlocksNotWatch: lastBlocksNotWatch,

		horizon:         horizon,
		addressProvider: addressProvider,
		builder:         builder,
		offchainHelper:  offchainHelper,
	}

	return result
}

func (s *Service) Run(ctx context.Context) {
	app.RunOverIncrementalTimer(ctx, s.log, s.serviceName, s.processNewBlocks, 5*time.Second, 5*time.Second)
}

func (s *Service) processNewBlocks(ctx context.Context) error {
	lastKnownBlock, err := s.offchainHelper.GetLastKnownBlockNumber()
	if err != nil {
		return errors.Wrap(err, "Failed to GetBlockCount")
	}

	// We only consider transactions, which have at least LastBlocksNotWatch + 1 confirmations
	lastBlockToProcess := lastKnownBlock - s.lastBlocksNotWatch

	if lastBlockToProcess <= s.lastProcessedBlock {
		// No new blocks to process
		return nil
	}

	for i := s.lastProcessedBlock + 1; i <= lastBlockToProcess; i++ {
		if app.IsCanceled(ctx) {
			return nil
		}

		err := s.processBlock(ctx, i)
		if err != nil {
			return errors.Wrap(err, "Failed to process Block", logan.F{"block_index": i})
		}

		if app.IsCanceled(ctx) {
			// Don't update lastProcessedBlock, because Block processing was probably not finished - ctx was canceled.
			return nil
		}

		s.lastProcessedBlock = i
	}

	return nil
}

func (s *Service) processBlock(ctx context.Context, blockIndex uint64) error {
	s.log.WithField("block_index", blockIndex).Debug("Processing block.")

	block, err := s.offchainHelper.GetBlock(blockIndex)
	if err != nil {
		return errors.Wrap(err, "Failed to get Block from BTCClient")
	}

	for _, tx := range block.TXs {
		if app.IsCanceled(ctx) {
			return nil
		}

		err := s.processTX(ctx, block.Hash, block.Timestamp, tx)
		if err != nil {
			return errors.Wrap(err, "Failed to process TX", logan.F{
				"block": block,
				"tx":    tx,
			})
		}
	}

	return nil
}

func (s *Service) processTX(ctx context.Context, blockHash string, blockTime time.Time, tx Tx) error {
	for i, out := range tx.Outs {
		if out.Address == "" {
			continue
		}

		accountAddress := s.addressProvider.AddressAt(ctx, blockTime, out.Address)
		if app.IsCanceled(ctx) {
			return nil
		}

		if accountAddress == nil {
			// This addr58 is not among our watch addresses - ignoring this TxOUT.
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
				"block_hash": blockHash,
				"block_time": blockTime,
				"tx_hash":    tx.Hash,
			})).WithField("min_deposit_amount_from_config", s.offchainHelper.GetMinDepositAmount()).
				Warn("Received deposit with too small amount.")
			continue
		}

		err := s.processDeposit(ctx, blockHash, blockTime, tx.Hash, i, out, *accountAddress)
		// TODO Retry processing Deposit - don't go to following Deposits.
		if err != nil {
			return errors.Wrap(err, "Failed to process Deposit", fields)
		}
	}

	return nil
}

func (s *Service) processDeposit(ctx context.Context, blockHash string, blockTime time.Time, txHash string, outIndex int, out Out, accountAddress string) error {
	s.log.WithFields(logan.F{
		"block_hash":      blockHash,
		"tx_hash":         txHash,
		"out_index":       outIndex,
		"offchain_addr":   out.Address,
		"offchain_amount": out.Value,
		"account_address": accountAddress,
	}).Debug("Processing deposit.")

	balanceID, err := s.getBalanceID(accountAddress)
	if err != nil {
		return errors.Wrap(err, "Failed to get BalanceID")
	}

	valueWithoutDepositFee := out.Value - s.offchainHelper.GetFixedDepositFee()
	emissionAmount := s.offchainHelper.ConvertToSystem(valueWithoutDepositFee)

	reference := s.offchainHelper.BuildReference(blockHash, txHash, out.Address, uint(outIndex), out.Value, 64)

	issuance := issuanceRequestOpt{
		Reference: reference,
		Receiver:  balanceID,
		Asset:     s.offchainHelper.GetAsset(),
		Amount:    emissionAmount,
		Details: resources.DepositDetails{
			TXHash: txHash,
			Price:  amount.One,
		}.Encode(),
	}

	err = s.processIssuance(ctx, blockHash, out.Address, accountAddress, issuance)
	if err != nil {
		return errors.Wrap(err, "Failed to send CoinEmissionRequest", logan.F{
			"balance_id":              balanceID,
			"converted_system_amount": emissionAmount,
			"reference":               reference,
		})
	}

	return nil
}

func (s *Service) processIssuance(ctx context.Context, blockHash, offchainAddress, accountAddress string, issuance issuanceRequestOpt) error {
	tx := s.craftIssuanceTX(issuance)

	// TODO Verify
	// TODO Remove after moving logic of the second signature to the verifier.
	tx = tx.Sign(s.additionalSigner)

	envelope, err := tx.Marshal()
	if err != nil {
		return errors.Wrap(err, "Failed to marshal TX into envelope")
	}

	logger := s.log.WithFields(logan.F{
		"block_hash":       blockHash,
		"offchain_address": offchainAddress,
		"account_address":  accountAddress,
		"issuance":         issuance,
	})

	result := s.horizon.Submitter().Submit(ctx, envelope)
	if result.Err != nil {
		// TODO Detect reference duplication errors and never log them
		// TODO Now any submit error is only logged and ignored - it's a problem.
		logger.WithFields(logan.F{
			"submit_response_raw":      string(result.RawResponse),
			"submit_response_tx_code":  result.TXCode,
			"submit_response_op_codes": result.OpCodes,
		}).WithError(result.Err).Error("Failed to submit CoinEmissionRequest Transaction to Horizon.")
		return nil
	}

	logger.Info("CoinEmissionRequest was sent successfully.")
	return nil
}

func (s *Service) getBalanceID(accountAddress string) (string, error) {
	balances, err := s.horizon.WithSigner(s.signer).Accounts().Balances(accountAddress)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get Account Balances from Horizon")
	}

	for _, b := range balances {
		if b.Asset == s.offchainHelper.GetAsset() {
			return b.BalanceID, nil
		}
	}

	// No BalanceID of the Offchain asset for the Account.
	return "", ErrNoBalanceID
}
