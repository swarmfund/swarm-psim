package btcsupervisor

import (
	"time"

	"context"
	"math/big"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2/types"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
	"gitlab.com/swarmfund/psim/psim/internal/resources"
	"gitlab.com/swarmfund/psim/psim/supervisor"
)

const (
	runnerName      = "btc_supervisor"
	baseAsset       = "SUN"
	referenceMaxLen = 64
)

var (
	lastProcessedBlock uint64
	errNilPrice        = errors.New("PriceAt of accountDataProvider returned nil Price.")
)

func (s *Service) processBTCBlocksInfinitely(ctx context.Context) {
	lastProcessedBlock = s.config.LastProcessedBlock
	app.RunOverIncrementalTimer(ctx, s.Log, runnerName, s.processNewBTCBlocks, 5*time.Second, 5*time.Second)
}

func (s *Service) processNewBTCBlocks(ctx context.Context) error {
	lastKnownBlock, err := s.btcClient.GetBlockCount()
	if err != nil {
		return errors.Wrap(err, "Failed to GetBlockCount")
	}

	// We only consider transactions, which have at least LastBlocksNotWatch + 1 confirmations
	lastBlockToProcess := lastKnownBlock - s.config.LastBlocksNotWatch

	if lastBlockToProcess <= lastProcessedBlock {
		// No new blocks to process
		return nil
	}

	for i := lastProcessedBlock + 1; i <= lastBlockToProcess; i++ {
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

		lastProcessedBlock = i
	}

	return nil
}

func (s *Service) processBlock(ctx context.Context, blockIndex uint64) error {
	s.Log.WithField("block_index", blockIndex).Debug("Processing block.")

	block, err := s.btcClient.GetBlock(blockIndex)
	if err != nil {
		return errors.Wrap(err, "Failed to get Block from BTCClient")
	}

	blockTime := block.MsgBlock().Header.Timestamp
	blockHash := block.Hash().String()

	for _, tx := range block.Transactions() {
		err := s.processTX(ctx, blockHash, blockTime, tx)
		if err != nil {
			// Tx hash is added into logs inside processTX.
			return errors.Wrap(err, "Failed to process TX", logan.F{"block_hash": blockHash})
		}
	}

	return nil
}

func (s *Service) processTX(ctx context.Context, blockHash string, blockTime time.Time, tx *btcutil.Tx) error {
	for i, out := range tx.MsgTx().TxOut {
		scriptClass, addrs, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, s.btcClient.GetNetParams())

		if scriptClass != txscript.PubKeyHashTy {
			// Output, which pays not to a pub-key-hash Address - just ignoring.
			// We only accept deposits to our Addresses which are all actually pay-to-pub-key-hash addresses.
			continue
		}

		addr58 := addrs[0].String()

		accountAddress := s.accountDataProvider.AddressAt(ctx, blockTime, addr58)
		if app.IsCanceled(ctx) {
			return nil
		}

		if accountAddress == nil {
			// This addr58 is not among our watch addresses - ignoring this TxOUT.
			continue
		}

		err = s.processDeposit(ctx, blockHash, blockTime, tx.Hash().String(), i, *out, addr58, *accountAddress)
		if err != nil {
			return errors.Wrap(err, "Failed to process deposit", logan.F{
				"block_hash":      blockHash,
				"block_time":      blockTime,
				"tx_hash":         tx.Hash().String(),
				"out_value":       out.Value,
				"out_index":       i,
				"btc_addr":        addr58,
				"account_address": accountAddress,
			})
		}
	}

	return nil
}

// TODO Check that amount is valid.
func (s *Service) processDeposit(ctx context.Context, blockHash string, blockTime time.Time, txHash string, outIndex int, out wire.TxOut, addr58, accountAddress string) error {
	s.Log.WithFields(logan.F{
		"block_hash":      blockHash,
		"tx_hash":         txHash,
		"btc_addr":        addr58,
		"account_address": accountAddress,
	}).
		Debug("Processing deposit.")

	price := s.accountDataProvider.PriceAt(ctx, blockTime)
	if price == nil {
		return errNilPrice
	}

	// TODO Check that amount is valid.
	// amount = value * price / 10^8
	div := new(big.Int).Mul(big.NewInt(100000000), big.NewInt(1))
	bigPrice := big.NewInt(*price)

	amount := new(big.Int).Mul(big.NewInt(int64(out.Value)), bigPrice)
	amount = amount.Div(amount, div)
	if !amount.IsUint64() {
		return errors.New("Amount value overflow.")
	}

	err := s.sendCoinEmissionRequest(blockHash, txHash, outIndex, accountAddress, amount.Uint64(), bigPrice.Int64())
	if err != nil {
		return errors.Wrap(err, "Failed to send CoinEmissionRequest", logan.F{
			"converted_amount": amount.Uint64(),
		})
	}

	return nil
}

func (s *Service) sendCoinEmissionRequest(blockHash, txHash string, outIndex int, accountAddress string, amount uint64, price int64) error {
	reference := bitcoin.BuildCoinEmissionRequestReference(txHash, outIndex)
	// Just in case. Reference must not be longer than 64.
	if len(reference) > referenceMaxLen {
		reference = reference[len(reference)-referenceMaxLen:]
	}

	// TODO Verify

	// TODO Verify
	//txBuilder := s.PrepareCERTx(reference, accountAddress, amount)
	//
	//envelope, err := txBuilder.Marshal64()
	//if err != nil {
	//	return errors.Wrap(err, "Failed to craft CoinEmissionRequests envelope")
	//}
	//
	//verifyPayload, err := json.Marshal(btcverify.VerifyRequest{
	//	Envelope:  *envelope,
	//	BlockHash: blockHash,
	//	TXHash:    txHash,
	//	OutIndex:  outIndex,
	//})
	//
	//s.SendCoinEmissionRequestForVerify(s.Ctx, verifyPayload)

	logger := s.Log.WithFields(logan.F{"account_address": accountAddress,
		"reference":  reference,
		"block_hash": blockHash,
		"tx_hash":    txHash,
		"out_index":  outIndex,
	})

	balances, err := s.horizon.WithSigner(s.config.Supervisor.SignerKP).Accounts().Balances(accountAddress)
	if err != nil {
		return errors.Wrap(err, "Failed to get Account from Horizon")
	}

	receiver := ""
	for _, b := range balances {
		if b.Asset == baseAsset {
			receiver = b.BalanceID
		}
	}

	// TODO Handle if no receiver.

	logger = logger.WithField("receiver", receiver)

	logger.Info("Sending CoinEmissionRequest.")

	tx := s.CraftIssuanceRequest(supervisor.IssuanceRequestOpt{
		Asset:     baseAsset,
		Reference: reference,
		Receiver:  receiver,
		Amount:    amount,
		Details: resources.DepositDetails{
			Source: txHash,
			Price:  types.Amount(price),
		}.Encode(),
	})

	// TODO Remove after moving logic of the second signature to the verifier.
	tx = tx.Sign(s.config.AdditionalSignerKP)

	envelope, err := tx.Marshal()

	result := s.horizon.Submitter().Submit(context.TODO(), envelope)
	if result.Err != nil {
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
