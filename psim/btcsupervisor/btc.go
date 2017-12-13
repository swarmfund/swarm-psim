package btcsupervisor

import (
	"encoding/json"

	"time"

	"github.com/piotrnar/gocoin/lib/btc"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
	"gitlab.com/swarmfund/psim/psim/btcverify"
)

const (
	// TODO Move to config
	lastBlocksNotWatch = 5 // We only consider transactions, which have 6+ confirmations
	runnerName         = "btc_supervisor"
)

var (
	lastProcessedBlock uint64
)

func (s *Service) processBTCBlocksInfinitely() {
	app.RunOverIncrementalTimer(s.Ctx, s.Log, runnerName, s.processNewBTCBlocks, 5*time.Second)
}

func (s *Service) processNewBTCBlocks() error {
	select {
	case <-s.Ctx.Done():
		return nil
	case <-s.addressQ.ReadinessWaiter():
		// Can continue
	}

	// addressQ is now ready - can proceed.
	lastKnownBlock, err := s.btcClient.GetBlockCount()
	if err != nil {
		return errors.Wrap(err, "Failed to GetBlockCount")
	}

	lastBlockToProcess := lastKnownBlock - lastBlocksNotWatch

	if lastBlockToProcess <= lastProcessedBlock {
		// No new blocks to process
		return nil
	}

	for i := lastProcessedBlock + 1; i <= lastBlockToProcess; i++ {
		if app.IsCanceled(s.Ctx) {
			return nil
		}

		err := s.processBlock(i)
		if err != nil {
			return errors.Wrap(err, "Failed to process Block", logan.Field("block_index", i))
		}

		if app.IsCanceled(s.Ctx) {
			// Don't update lastProcessedBlock, because Block processing was probably not finished - ctx was canceled.
			return nil
		}

		lastProcessedBlock = i
	}

	return nil
}

func (s *Service) processBlock(blockIndex uint64) error {
	s.Log.WithField("block_index", blockIndex).Debug("Processing block")

	block, err := s.btcClient.GetBlock(blockIndex)
	if err != nil {
		return errors.Wrap(err, "Failed to get Block from BTCClient")
	}

	for _, tx := range block.Txs {
		err := s.processTX(block.Hash.String(), *tx)
		if err != nil {
			return errors.Wrap(err, "Failed to process TX", logan.Field("tx_hash", tx.Hash.String()))
		}
	}

	return nil
}

func (s *Service) processTX(blockHash string, tx btc.Tx) error {
	for i, out := range tx.TxOut {
		addr := btc.NewAddrFromPkScript(out.Pk_script, s.btcClient.IsTestnet())
		if addr == nil {
			// Somebody is playing with sending BTC not to an address - just ignore this
			continue
		}

		addr58 := addr.String()

		accountID := s.addressQ.GetAccountID(addr58)
		if accountID == "" {
			// This addr58 is not among our watch addresses - ignoring this TxOUT
			continue
		}

		// TODO Make sure amount is valid (maybe need to * or / by 10^N)
		err := s.sendCoinEmissionRequest(blockHash, tx.Hash.String(), i, accountID, out.Value)
		if err != nil {
			errors.Wrap(err, "Failed to send CoinEmissionRequest",
				logan.Field("tx", tx).Add("out", out).Add("account_id", accountID))
		}
	}

	return nil
}

func (s *Service) sendCoinEmissionRequest(blockHash, txHash string, outIndex int, receiver string, amount uint64) error {
	reference := bitcoin.BuildCoinEmissionRequestReference(txHash, outIndex)

	cerExists, err := s.CheckCoinEmissionRequestExistence(reference)
	if err != nil {
		return errors.Wrap(err, "Failed to check CoinEmissionRequest existence")
	}

	if cerExists {
		return nil
	}

	s.Log.WithField("reference", reference).WithField("receiver", receiver).Info("Sending CoinEmissionRequest")

	// TODO
	envelope, err := s.PrepareIREnvelope(reference, receiver, bitcoin.Asset, amount)
	if err != nil {
		return errors.Wrap(err, "Failed to craft CoinEmissionRequests envelope")
	}

	verifyPayload, err := json.Marshal(btcverify.VerifyRequest{
		Envelope:  *envelope,
		BlockHash: blockHash,
		TXHash:    txHash,
		OutIndex:  outIndex,
	})

	s.SendCoinEmissionRequest(s.Ctx, verifyPayload)

	return nil
}
