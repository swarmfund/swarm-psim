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
)

var (
	lastProcessedBlock uint64 = 1255032
)

// TODO runner func must receive ctx
func (s *Service) processBTCBlocksInfinitely() {
	// TODO runner func must receive ctx
	ctx := s.Ctx
	app.RunOverIncrementalTimer(ctx, s.Log, runnerName, s.processNewBTCBlocks, 5*time.Second)
}

func (s *Service) processNewBTCBlocks(ctx context.Context) error {
	// addressProvider is now ready - can proceed.
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

	for _, tx := range block.Txs {
		blockHash := block.Hash.String()

		// TODO Check that block.MedianPastTime is time of Block.
		err := s.processTX(ctx, blockHash, time.Now().UTC(), *tx)
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

		accountAddress := s.addressProvider.AddressAt(ctx, blockTime, addr58)
		if app.IsCanceled(ctx) {
			return nil
		}

		if accountAddress == nil {
			// This addr58 is not among our watch addresses - ignoring this TxOUT
			continue
		}

		price := s.addressProvider.PriceAt(s.Ctx, blockTime)
		if price == nil {
			return errors.New("price not set")
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

	//cerExists, err := s.CheckCoinEmissionRequestExistence(reference)
	//if err != nil {
	//	return errors.Wrap(err, "Failed to check CoinEmissionRequest existence")
	//}
	//
	//if cerExists {
	//	return nil
	//}

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

	fields := logan.Field("account_address", accountAddress).Add("reference", reference).
		Add("block_hash", blockHash).Add("tx_hash", txHash).Add("out_index", outIndex)

	balanceID, err := s.addressProvider.BalanceID(ctx, accountAddress)
	if err != nil {
		return errors.Wrap(err, "Failed to get BalanceID by AccountAddress", logan.Field("reference", reference))
	}
	if balanceID == nil {
		s.Log.WithFields(fields).Error("BalanceID is empty.")
		return nil
	}

	s.Log.WithFields(fields).Info("Sending CoinEmissionRequest.")

	err = s.PrepareCERTx(reference, *balanceID, amount).Submit()
	if err != nil {
		return errors.Wrap(err, "Failed to submit CoinEmissionRequest", logan.Field("reference", reference).
			Add("balance_id", balanceID))
	}

	s.Log.WithFields(fields).Info("CoinEmissionRequest was sent successfully.")
	return nil
}
