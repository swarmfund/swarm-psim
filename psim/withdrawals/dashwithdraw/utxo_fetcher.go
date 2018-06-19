package dashwithdraw

import (
	"context"
	"time"

	"sync"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
)

type UTXO struct {
	bitcoin.Out
	Value      int64
	IsInactive bool
}

func (h *CommonDashHelper) fetchUTXOsInfinitely(ctx context.Context, blockToStart uint64) {
	h.log.WithField("start_block", blockToStart).Info("Starting UTXOs fetcher.")

	blockStream := make(chan *btcutil.Block, 5)
	var lastKnownBlock uint64

	running.UntilSuccess(ctx, h.log, "last_known_block_getter", func(ctx context.Context) (bool, error) {
		var err error
		lastKnownBlock, err = h.btcClient.GetBlockCount(ctx)
		if err != nil {
			return false, errors.Wrap(err, "Failed to GetBlockCount (last known Block)")
		}

		return true, nil
	}, 5*time.Second, 5*time.Minute)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		h.streamOffchainBlocks(ctx, blockToStart, lastKnownBlock, blockStream)

		wg.Done()
	}()

	for {
		if running.IsCancelled(ctx) {
			break // the for
		}

		select {
		case <-ctx.Done():
			break // the select
		case block := <-blockStream:
			if running.IsCancelled(ctx) {
				break // the select
			}

			running.UntilSuccess(ctx, h.log, "block_utxo_processor", func(ctx context.Context) (bool, error) {
				err := h.processBlockUTXOs(block)
				if err != nil {
					return false, err
				}

				return true, nil
			}, 5*time.Second, time.Hour)
		}
	}

	wg.Wait()
}

func (h *CommonDashHelper) streamOffchainBlocks(ctx context.Context, blockToStart, lastKnownBlock uint64, blockStream chan<- *btcutil.Block) {
	lastProcessedBlock := blockToStart
	initialCtx, finishedInitial := context.WithCancel(ctx)

	running.WithBackOff(initialCtx, h.log, "initial_blocks_streamer", func(ctx context.Context) error {
		if lastKnownBlock <= lastProcessedBlock {
			// No new (still unprocessed) Blocks.
			h.log.Info("All existing Offchain Blocks(and UTXOs) are fetched.")
			close(h.utxoFetched)
			finishedInitial()
			return nil
		}

		blockNumber := lastProcessedBlock + 1

		block, err := h.btcClient.GetBlock(blockNumber)
		if err != nil {
			return errors.Wrap(err, "Failed to GetBlock", logan.F{"block_number": blockNumber})
		}
		block.SetHeight(int32(blockNumber))

		select {
		case <-ctx.Done():
			return nil
		case blockStream <- block:
			lastProcessedBlock = blockNumber
		}

		return nil
	}, 0, 5*time.Second, time.Hour)

	running.WithBackOff(ctx, h.log, "blocks_streamer", func(ctx context.Context) error {
		lastKnownBlock, err := h.btcClient.GetBlockCount(ctx)
		if err != nil {
			return errors.Wrap(err, "Failed to GetBlockCount (last known Block)")
		}

		if lastKnownBlock <= lastProcessedBlock {
			// No new (still unprocessed) Blocks.
			return nil
		}

		blockNumber := lastProcessedBlock + 1

		block, err := h.btcClient.GetBlock(blockNumber)
		if err != nil {
			return errors.Wrap(err, "Failed to GetBlock", logan.F{"block_number": blockNumber})
		}
		block.SetHeight(int32(blockNumber))

		select {
		case <-ctx.Done():
			return nil
		case blockStream <- block:
			lastProcessedBlock = blockNumber
		}

		return nil
	}, 10*time.Second, 5*time.Second, time.Hour)
}

func (h *CommonDashHelper) processBlockUTXOs(block *btcutil.Block) error {
	h.log.WithField("block_number", block.Height()).Info("Processing Block UTXOs.")

	for _, tx := range block.Transactions() {
		for i, out := range tx.MsgTx().TxOut {
			h.saveUTXOIfOurs(*out, tx.Hash().String(), i)
		}

		for _, in := range tx.MsgTx().TxIn {
			wasRemoved := h.coinSelector.TryRemoveUTXO(bitcoin.Out{
				TXHash: in.PreviousOutPoint.Hash.String(),
				Vout:   in.PreviousOutPoint.Index,
			})

			if wasRemoved {
				h.log.WithFields(logan.F{
					"tx_hash":    tx.Hash().String(),
					"out_number": in.PreviousOutPoint.Index,
				}).Debug("Found our UTXO being spent - removing it locally.")
			}
		}
	}

	return nil
}

func (h *CommonDashHelper) saveUTXOIfOurs(out wire.TxOut, txHash string, outNumber int) {
	logger := h.log.WithFields(logan.F{
		"tx_hash":    txHash,
		"out_number": outNumber,
	})

	scriptClass, addrs, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, h.netParams)
	if err != nil {
		// Just a weird TX Output in the Blockchain - just ignoring.
		return
	}

	if scriptClass != txscript.PubKeyHashTy {
		// Output, which pays not to a pub-key-hash Address - just ignoring.
		// We only funnel BTC from our Addresses, which are all actually pay-to-pub-key-hash Addresses.
		return
	}

	if len(addrs) == 0 {
		logger.Error("Found Output with empty Addresses parsed from PKScript.")
		return
	}
	addr58 := addrs[0].String()
	logger = logger.WithField("addr", addr58)

	if addr58 != h.config.HotWalletAddress {
		// This pay-to-pub-key-hash Address is not ours.
		return
	}

	// Found our TX Output.
	utxo := UTXO{
		Out: bitcoin.Out{
			TXHash: txHash,
			Vout:   uint32(outNumber),
		},
		Value: out.Value,
	}
	logger.WithField("utxo", utxo).Debug("Found our UTXO.")

	h.coinSelector.AddUTXO(utxo)
	return
}
