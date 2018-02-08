package btcfunnel

import (
	"context"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
)

// FetchExistingBlocks fetches old Blocks and streams the Outputs in
// the TXs of these Blocks which belong to our Addresses (see `s.addrToPriv`)
// into `s.outCh`.
func (s *Service) fetchExistingBlocks(ctx context.Context) error {
	s.log.Info("Started fetching existing Blocks.")

	lastExistingBlock, err := s.btcClient.GetBlockCount()
	if err != nil {
		return errors.Wrap(err, "Failed to GetBlockCount")
	}

	for s.lastProcessedBlock < lastExistingBlock {
		if app.IsCanceled(ctx) {
			return nil
		}

		blockNumber := s.lastProcessedBlock + 1
		s.log.WithField("block_number", blockNumber).Info("Processing existing Block.")

		err := s.processBlock(blockNumber)
		if err != nil {
			return errors.Wrap(err, "Failed to process Block")
		}
	}

	s.log.WithField("last_processed_block", s.lastProcessedBlock).Info("Finished fetching existing Blocks.")
	close(s.outCh)
	return nil
}

func (s *Service) processBlock(blockNumber uint64) error {
	block, err := s.btcClient.GetBlock(blockNumber)
	if err != nil {
		return errors.Wrap(err, "Failed to GetBlock", logan.F{"block_number": blockNumber})
	}

	err = s.streamOurOuts(block)
	if err != nil {
		return errors.Wrap(err, "Failed to stream our Outputs of the Block", logan.F{
			"block_number": blockNumber,
			"block":        block,
		})
	}

	// Block was successfully processed
	s.lastProcessedBlock = blockNumber
	return nil
}

func (s *Service) streamOurOuts(block *btcutil.Block) error {
	for _, tx := range block.Transactions() {
		fields := logan.F{
			"tx_hash": tx.Hash().String(),
		}

		for i, out := range tx.MsgTx().TxOut {
			fields["out_number"] = i

			scriptClass, addrs, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, s.btcClient.GetNetParams())
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
				s.log.WithFields(fields).Error("Found Output with empty Addresses parsed from PKScript.")
				continue
			}
			addr58 := addrs[0].String()
			fields["addr"] = addr58

			_, ok := s.addrToPriv[addr58]
			if !ok {
				// This pay-to-pub-key-hash Address is not ours.
				continue
			}

			// Found our TX Output.
			s.outCh <- bitcoin.Out{
				tx.Hash().String(),
				uint32(i),
			}
		}
	}

	return nil
}

// TODO SYNCHRONIZE
func (s *Service) fetchNewBlock(ctx context.Context) error {
	lastKnownBlock, err := s.btcClient.GetBlockCount()
	if err != nil {
		return errors.Wrap(err, "Failed to GetBlockCount")
	}

	if lastKnownBlock <= s.lastProcessedBlock {
		// No new (still unprocessed) Blocks.
		return nil
	}

	// TODO SYNCHRONIZE
	s.outCh = make(chan bitcoin.Out, 1000)
	go s.listenOutStream(ctx)

	blockNumber := s.lastProcessedBlock + 1
	s.log.WithField("block_number", blockNumber).Info("Processing newly appeared Block.")

	err = s.processBlock(blockNumber)
	if err != nil {
		return errors.Wrap(err, "Failed to process Block")
	}

	close(s.outCh)
	return nil
}
