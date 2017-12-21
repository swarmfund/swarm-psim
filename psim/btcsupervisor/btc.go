package btcsupervisor

import (
	"time"

	"context"
	"math/big"

	"github.com/piotrnar/gocoin/lib/btc"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
)

const (
	runnerName = "btc_supervisor"
	baseAsset  = "SUN"
)

var (
	// FIXME Use config
	lastProcessedBlock uint64 = 1255110
)

func (s *Service) processBTCBlocksInfinitely(ctx context.Context) {
	lastProcessedBlock = s.config.LastProcessedBlock
	app.RunOverIncrementalTimer(ctx, s.Log, runnerName, s.processNewBTCBlocks, 5*time.Second)
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
			return errors.Wrap(err, "Failed to process Block", logan.Field("block_index", i))
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

	blockTime := time.Unix(int64(block.BlockTime()), 0)

	for _, tx := range block.Txs {
		blockHash := block.Hash.String()

		err := s.processTX(ctx, blockHash, blockTime, *tx)
		if err != nil {
			// Tx hash is added into logs inside processTX.
			return errors.Wrap(err, "Failed to process TX", logan.Field("block_hash", blockHash))
		}
	}

	return nil
}

func (s *Service) processTX(ctx context.Context, blockHash string, blockTime time.Time, tx btc.Tx) error {
	for i, out := range tx.TxOut {
		addr := btc.NewAddrFromPkScript(out.Pk_script, s.btcClient.IsTestnet())
		if addr == nil {
			// Somebody is playing with sending BTC not to an address - just ignore this
			continue
		}

		addr58 := addr.String()

		accountAddress := s.accountDataProvider.AddressAt(ctx, blockTime, addr58)
		if app.IsCanceled(ctx) {
			return nil
		}

		if accountAddress == nil {
			// This addr58 is not among our watch addresses - ignoring this TxOUT
			continue
		}

		s.Log.WithField("block_hash", blockHash).WithField("btc_addr", addr58).
			WithField("account_address", accountAddress).Debug("Found our watch BTC Address.")

		price := s.accountDataProvider.PriceAt(ctx, blockTime)
		if price == nil {
			return errors.From(errors.New("PriceAt of accountDataProvider returned nil price."),
				logan.Field("block_hash", blockHash).Add("btc_addr", addr58).Add("account_address", accountAddress).
					Add("block_time", blockTime))
		}

		// amount = value * price / 10^8
		div := new(big.Int).Mul(big.NewInt(100000000), big.NewInt(1))
		bigPrice := big.NewInt(*price)

		amount := new(big.Int).Mul(big.NewInt(int64(out.Value)), bigPrice)
		amount = amount.Div(amount, div)
		if !amount.IsUint64() {
			return errors.New("value overflow")
		}

		// Don't take hash outside of for loop, as it's will be needed not more than once during whole processTX(), usually will not be used at all.
		txHash := tx.Hash.String()
		err := s.sendCoinEmissionRequest(ctx, blockHash, txHash, i, *accountAddress, amount.Uint64())
		if err != nil {
			return errors.Wrap(err, "Failed to send CoinEmissionRequest",
				logan.Field("account_address", *accountAddress).
					Add("tx_hash", txHash).
					Add("out_index", out.VoutCount).
					Add("out_value", out.Value).
					Add("out_addr_58", addr58).
					Add("converted_amount", amount.Uint64()),
			)
		}
	}

	return nil
}

func (s *Service) sendCoinEmissionRequest(ctx context.Context, blockHash, txHash string, outIndex int, accountAddress string, amount uint64) error {
	reference := bitcoin.BuildCoinEmissionRequestReference(txHash, outIndex)
	// Just in case. Reference must not be longer than 64.
	if len(reference) > 64 {
		reference = reference[len(reference)-64:]
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

	fields := logan.F{"account_address": accountAddress,
						"reference": reference,
						"block_hash": blockHash,
						"tx_hash": txHash,
						"out_index": outIndex,}

	// TODO Move getting of BalanceID into accountDataProvider.
	account, err := s.horizon.AccountSigned(s.config.Supervisor.SignerKP, accountAddress)
	if err != nil {
		return err
	}

	if account == nil {
		return errors.New("Horizon returned nil Account.")
	}

	receiver := ""
	for _, b := range account.Balances {
		if b.Asset == baseAsset {
			receiver = b.BalanceID
		}
	}

	// TODO Handle if no receiver.

	fields = fields.Add("receiver", receiver)

	s.Log.WithFields(fields).Info("Sending CoinEmissionRequest.")

	err = s.PrepareCERTx(reference, receiver, amount).
		Submit()
	if err != nil {
		s.Log.WithError(err).Error("Failed to submit CoinEmissionRequest Transaction.")
		return nil
	}

	s.Log.WithFields(fields).Info("CoinEmissionRequest was sent successfully.")
	return nil
}
