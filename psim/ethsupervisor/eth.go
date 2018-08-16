package ethsupervisor

import (
	"math/big"
	"time"

	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/ethsupervisor/internal"
	"gitlab.com/swarmfund/psim/psim/internal/resources"
	"gitlab.com/swarmfund/psim/psim/supervisor"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/regources"
)

// TODO defer
func (s *Service) processBlocks(ctx context.Context) {
	// TODO Listen to both s.blockCh and ctx.Done() in select.
	for blockNumber := range s.blocksCh {
		if app.IsCanceled(ctx) {
			return
		}

		entry := s.Log.WithField("block_number", blockNumber)
		entry.Debug("Processing block.")

		block, err := s.eth.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
		if err != nil {
			entry.Error("Failed to get block.")
			s.blocksCh <- blockNumber
			continue
		}

		if block == nil {
			entry.Error("Got nil block from node.")
			continue
		}

		for _, tx := range block.Transactions() {
			if app.IsCanceled(ctx) {
				return
			}

			if tx == nil {
				continue
			}

			s.txCh <- internal.Transaction{
				Timestamp:   time.Unix(block.Time().Int64(), 0),
				BlockNumber: block.NumberU64(),
				Transaction: *tx,
			}
		}
	}
}

// TODO eth.SubscribeNewHead
func (s *Service) watchHeight(ctx context.Context) {
	cursor := new(big.Int).Set(s.config.Cursor)

	go func() {
		ticker := time.NewTicker(10 * time.Second)

		// TODO Listen to both ticker.C and ctx.Done() in select.
		for ; ; <-ticker.C {
			if app.IsCanceled(ctx) {
				return
			}

			head, err := s.eth.BlockByNumber(ctx, nil)
			if err != nil {
				s.Log.WithError(err).Error("failed to get block count")
				continue
			}

			s.Log.WithField("height", head.NumberU64()).Debug("fetched new head")

			for head.NumberU64()-s.config.Confirmations > cursor.Uint64() {
				if app.IsCanceled(ctx) {
					return
				}

				s.blocksCh <- cursor.Uint64()
				cursor.Add(cursor, big.NewInt(1))
			}
		}
	}()
}

func (s *Service) processTXs(ctx context.Context) {
	for tx := range s.txCh {
		if app.IsCanceled(ctx) {
			return
		}

		entry := s.Log.WithField("tx", tx.Hash().Hex())

		for {
			// TODO Incremental!
			if err := s.processTX(ctx, tx); err != nil {
				entry.WithError(err).Error("Failed to process TX.")
				continue
			}
			break
		}
	}
}

// TODO Refactor me.
func (s *Service) processTX(ctx context.Context, tx internal.Transaction) (err error) {
	defer func() {
		if rvr := recover(); rvr != nil {
			err = errors.FromPanic(rvr)
		}
	}()

	// tx amount exceeds deposit threshold
	if tx.Value().Cmp(s.depositThreshold) == -1 {
		return nil
	}

	// tx has destination
	if tx.To() == nil {
		return nil
	}

	// address is watched
	address := s.state.ExternalAccountAt(ctx, tx.Timestamp, s.config.ExternalSystem, tx.To().Hex())
	if app.IsCanceled(ctx) {
		return nil
	}

	if address == nil {
		return nil
	}

	entry := s.Log.WithFields(logan.F{
		"tx_hash":     tx.Hash().Hex(),
		"eth_address": tx.To().String(),
	})
	entry.Info("Found deposit.")

	receiver := s.state.Balance(ctx, *address, s.config.DepositAsset)
	if app.IsCanceled(ctx) {
		return nil
	}

	if receiver == nil {
		entry.Error("balance not found, skipping tx")
		return nil
	}

	// amount = value * price / 10^18
	ethPrecision := new(big.Int).Mul(big.NewInt(1000000000), big.NewInt(1000000000))
	valueWithoutDepositFee := tx.Value().Sub(tx.Value(), s.config.FixedDepositFee)
	emissionAmount := new(big.Int).Mul(valueWithoutDepositFee, big.NewInt(amount.One))
	emissionAmount = emissionAmount.Div(emissionAmount, ethPrecision)
	if !emissionAmount.IsUint64() {
		entry.Error("amount overflow, skipping tx")
		return nil
	}

	reference := tx.Hash().Hex()
	// yoba eth hex trimming
	if len(reference) > 64 {
		reference = reference[len(reference)-64:]
	}
	request := s.CraftIssuanceRequest(supervisor.IssuanceRequestOpt{
		Asset:     s.config.DepositAsset,
		Reference: reference,
		Receiver:  *receiver,
		Amount:    emissionAmount.Uint64(),
		Details: regources.DepositDetails{
			TXHash: tx.Hash().Hex(),
		}.Encode(),
	})

	entry = entry.WithFields(logan.F{
		"emission_amount": emissionAmount.Uint64(),
		"reference":       reference,
		"receiver":        *receiver,
	})

	envelope, err := request.Marshal()
	if err != nil {
		return errors.Wrap(err, "Failed to marshal IssuanceRequest Envelope")
	}

	if result := s.Horizon.Submitter().Submit(ctx, envelope); result.Err != nil {
		entry := entry.WithFields(logan.F{
			"block":                    tx.BlockNumber,
			"submit_response_raw":      string(result.RawResponse),
			"submit_response_tx_code":  result.TXCode,
			"submit_response_op_codes": result.OpCodes,
		}).WithError(result.Err)

		if len(result.OpCodes) == 1 {
			switch result.OpCodes[0] {
			// safe to move on
			case "op_reference_duplication":
				entry.Info("Met op_reference_duplication in response from Horizon - skipping.")
				return nil
			}
		}

		entry.Error("Failed to submit CoinEmissionRequest Transaction.")
		return err
	}

	entry.Info("issuance request submitted")

	return nil
}
