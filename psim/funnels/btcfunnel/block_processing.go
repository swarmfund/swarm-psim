package btcfunnel

import (
	"context"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
)

// FetchExistingBlocks fetches old Blocks and streams the Outputs of
// the TXs of these Blocks which belong to our Addresses (see `s.addrToPriv`)
// into `s.outsCh`.
func (s *Service) fetchExistingBlocks(ctx context.Context) (bool, error) {
	s.log.Info("Started fetching existing Blocks.")
	lastExistingBlock, err := s.btcClient.GetBlockCount(ctx)
	if err != nil {
		return false, errors.Wrap(err, "Failed to GetBlockCount")
	}
	if running.IsCancelled(ctx) {
		return false, nil
	}

	outsCh := make(chan bitcoin.Out, outChanSize)
	defer close(outsCh)
	go s.listenOurOutsStream(ctx, outsCh)

	var totalOurOutsStreamed int
	for s.lastProcessedBlock < lastExistingBlock {
		if running.IsCancelled(ctx) {
			// outCh is being closed in defer
			return true, nil
		}

		blockNumber := s.lastProcessedBlock + 1
		s.log.WithField("block_number", blockNumber).Info("Processing existing Block.")

		ourOutsStreamed, err := s.processBlock(blockNumber, outsCh)
		if err != nil {
			// outCh is being closed in defer
			return false, errors.Wrap(err, "Failed to process Block", logan.F{"block_number": blockNumber})
		}

		totalOurOutsStreamed += ourOutsStreamed
	}

	s.log.WithFields(logan.F{
		"last_processed_block": s.lastProcessedBlock,
		"our_outputs_found":    totalOurOutsStreamed,
	}).Info("Finished fetching existing Blocks.")

	// outCh is being closed in defer
	return true, nil
}

func (s *Service) fetchNewBlock(ctx context.Context) error {
	lastKnownBlock, err := s.btcClient.GetBlockCount(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to GetBlockCount")
	}
	if running.IsCancelled(ctx) {
		return nil
	}

	if lastKnownBlock <= s.lastProcessedBlock {
		// No new (still unprocessed) Blocks.
		return nil
	}

	outsCh := make(chan bitcoin.Out, outChanSize)
	go s.listenOurOutsStream(ctx, outsCh)

	blockNumber := s.lastProcessedBlock + 1
	s.log.WithField("block_number", blockNumber).Info("Processing newly appeared Block.")

	_, err = s.processBlock(blockNumber, outsCh)
	if err != nil {
		close(outsCh)
		return errors.Wrap(err, "Failed to process Block", logan.F{"block_number": blockNumber})
	}

	close(outsCh)
	return nil
}

func (s *Service) processBlock(blockNumber uint64, outsCh chan<- bitcoin.Out) (ourOutsStreamed int, err error) {
	block, err := s.btcClient.GetBlock(blockNumber)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to GetBlock")
	}

	ourOutsStreamed = s.streamBlockOutputs(block, outsCh)
	s.lastProcessedBlock = blockNumber
	return ourOutsStreamed, nil
}

func (s *Service) streamBlockOutputs(block *btcutil.Block, outsCh chan<- bitcoin.Out) (totalStreamed int) {
	var totalOurOutsStreamed int

	for _, tx := range block.Transactions() {
		outsStreamed := s.streamTXOutputs(tx, outsCh)
		totalOurOutsStreamed += outsStreamed
	}

	return totalOurOutsStreamed
}

func (s *Service) streamTXOutputs(tx *btcutil.Tx, outsCh chan<- bitcoin.Out) int {
	var totalStreamed int
	logger := s.log.WithField("tx_hash", tx.Hash().String())

	for i, out := range tx.MsgTx().TxOut {
		logger = logger.WithField("out_number", i)

		scriptClass, addrs, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, s.netParams)
		if err != nil {
			// Just a weird TX Output in the Blockchain - just ignoring.
			continue
		}

		if scriptClass != txscript.PubKeyHashTy {
			// Output, which pays not to a pub-key-hash Address - just ignoring.
			// We only funnel BTC from our Addresses, which are all actually pay-to-pub-key-hash Addresses.
			continue
		}

		if len(addrs) == 0 {
			logger.Error("Found Output with empty Addresses parsed from PKScript.")
			continue
		}
		addr58 := addrs[0].String()
		logger = logger.WithField("addr", addr58)

		_, ok := s.addrToPriv[addr58]
		if !ok {
			// This pay-to-pub-key-hash Address is not ours.
			continue
		}

		// Found our TX Output.
		outsCh <- bitcoin.Out{
			tx.Hash().String(),
			uint32(i),
		}

		totalStreamed += 1
	}

	return totalStreamed
}
