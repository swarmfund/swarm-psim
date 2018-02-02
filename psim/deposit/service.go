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

// AddressProvider must be implemented by WatchAddress storage to pass into Service constructor.
type AddressProvider interface {
	AddressAt(ctx context.Context, t time.Time, btcAddress string) (tokendAddress *string)
}

// TODO Comment
type OffchainHelper interface {
	// TODO Method for getting errors from addrstate chan

	//AddressAt(ctx context.Context, t time.Time, btcAddress string) (tokendAddress *string)

	// TODO Comments
	GetLastKnownBlockNumber() (uint64, error)
	GetBlock(number uint64) (*Block, error)
	GetMinDepositAmount() uint64
	GetFixedDepositFee() uint64
	ConvertToSystem(offchainAmount uint64) (systemAmount uint64)
	GetAsset() string
	BuildReference(blockHash, txHash, offchainAddress string, outIndex uint, amount uint64, maxLen int) string
}

// Service implements app.Service interface, it supervises Offchain blockchain
// and send CoinEmissionRequests to Horizon if arrived deposit detected.
type Service struct {
	log    *logan.Entry
	source keypair.Address
	signer keypair.Full
	// TODO Remove after moving logic of the second signature to the verifier.
	additionalSigner   keypair.Full
	serviceName        string
	lastProcessedBlock uint64
	lastBlocksNotWatch uint64

	// TODO Interface
	horizon         *horizon.Connector
	builder         *xdrbuild.Builder
	offchainHelper  OffchainHelper
	addressProvider AddressProvider
}

// New is constructor for the btcsupervisor Service.
func New(
	log *logan.Entry,
	source keypair.Address,
	signer keypair.Full,
	additionalSigner keypair.Full,
	serviceName string,
	lastProcessedBlock,
	lastBlocksNotWatch uint64,
	horizon *horizon.Connector,
	builder *xdrbuild.Builder,
	offchainHelper OffchainHelper,
	addressProvider AddressProvider) *Service {

	result := &Service{
		log:                log,
		source:             source,
		signer:             signer,
		additionalSigner:   additionalSigner,
		serviceName:        serviceName,
		lastProcessedBlock: lastProcessedBlock,
		lastBlocksNotWatch: lastBlocksNotWatch,

		horizon:         horizon,
		builder:         builder,
		offchainHelper:  offchainHelper,
		addressProvider: addressProvider,
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
		err := s.processTX(ctx, block.Hash, block.Timestamp, tx)
		if err != nil {
			// Tx hash is added into logs inside processTX.
			return errors.Wrap(err, "Failed to process TX", logan.F{
				"block":  block,
				"tx_hex": tx,
			})
		}
	}

	return nil
}

func (s *Service) processTX(ctx context.Context, blockHash string, blockTime time.Time, tx Tx) error {
	for i, out := range tx.Outs {
		accountAddress := s.addressProvider.AddressAt(ctx, blockTime, out.Address)
		if app.IsCanceled(ctx) {
			return nil
		}

		if accountAddress == nil {
			// This addr58 is not among our watch addresses - ignoring this TxOUT.
			continue
		}

		fields := logan.F{
			"block_hash":      blockHash,
			"block_time":      blockTime,
			"tx_hash":         tx.Hash,
			"out_value":       out.Value,
			"out_index":       i,
			"offchain_addr":   out.Address,
			"account_address": accountAddress,
		}

		if out.Value < s.offchainHelper.GetMinDepositAmount() {
			s.log.WithFields(fields).WithField("min_deposit_amount_from_config", s.offchainHelper.GetMinDepositAmount()).
				Warn("Received deposit with too small amount.")
			continue
		}

		err := s.processDeposit(ctx, blockHash, blockTime, tx.Hash, i, out, *accountAddress)
		if err != nil {
			return errors.Wrap(err, "Failed to process deposit", fields)
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

	valueWithoutDepositFee := out.Value - s.offchainHelper.GetFixedDepositFee()
	emissionAmount := s.offchainHelper.ConvertToSystem(valueWithoutDepositFee)

	reference := s.offchainHelper.BuildReference(blockHash, txHash, out.Address, uint(outIndex), out.Value, 64)

	err := s.sendCoinEmissionRequest(blockHash, txHash, out.Address, accountAddress, reference, emissionAmount)
	if err != nil {
		return errors.Wrap(err, "Failed to send CoinEmissionRequest", logan.F{
			"converted_system_amount": emissionAmount,
		})
	}

	return nil
}

func (s *Service) sendCoinEmissionRequest(blockHash, txHash, offchainAddress, accountAddress, reference string, emissionAmount uint64) error {
	// TODO Verify

	logger := s.log.WithFields(logan.F{
		"block_hash":       blockHash,
		"tx_hash":          txHash,
		"offchain_address": offchainAddress,
		"account_address":  accountAddress,
		"reference":        reference,
	})

	// TODO Move getting BalanceID to separate method
	balances, err := s.horizon.WithSigner(s.signer).Accounts().Balances(accountAddress)
	if err != nil {
		return errors.Wrap(err, "Failed to get Account Balances from Horizon")
	}

	receiver := ""
	for _, b := range balances {
		if b.Asset == s.offchainHelper.GetAsset() {
			receiver = b.BalanceID
		}
	}

	// TODO Handle if no receiver.

	logger = logger.WithField("receiver", receiver)

	tx := s.craftIssuanceTX(issuanceRequestOpt{
		Asset:     s.offchainHelper.GetAsset(),
		Reference: reference,
		Receiver:  receiver,
		Amount:    emissionAmount,
		Details: resources.DepositDetails{
			Source: txHash,
			Price:  amount.One,
		}.Encode(),
	})

	// TODO Remove after moving logic of the second signature to the verifier.
	tx = tx.Sign(s.additionalSigner)

	envelope, err := tx.Marshal()
	if err != nil {
		return errors.Wrap(err, "Failed to marshal TX into envelope")
	}

	result := s.horizon.Submitter().Submit(context.TODO(), envelope)
	if result.Err != nil {
		// TODO Detect reference duplication errors and never log them
		// TODO Now any submit error is only logged and ignored - it's a problem.
		logger.WithFields(logan.F{
			"submit_response_raw":      string(result.RawResponse),
			"submit_response_tx_code":  result.TXCode,
			"submit_response_op_codes": result.OpCodes,
		}).WithError(result.Err).Error("Failed to submit CoinEmissionRequest Transaction.")
		return nil
	}

	logger.Info("CoinEmissionRequest was sent successfully.")
	return nil
}
