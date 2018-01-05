package ethsupervisor

import (
	"math/big"
	"time"

	"context"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/horizon-connector/v2/types"
	"gitlab.com/swarmfund/psim/psim/ethsupervisor/internal"
	"gitlab.com/swarmfund/psim/psim/internal/resources"
	"gitlab.com/swarmfund/psim/psim/supervisor"
)

// TODO defer
func (s *Service) processBlocks(ctx context.Context) {
	for blockNumber := range s.blocksCh {
		s.Log.WithField("number", blockNumber).Debug("processing block")
		block, err := s.eth.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
		if err != nil {
			s.Log.WithField("block_number", blockNumber).Error("failed to get block")
			s.blocksCh <- blockNumber
			continue
		}

		if block == nil {
			s.Log.WithField("block_number", blockNumber).Error("missing block")
			continue
		}

		for _, tx := range block.Transactions() {
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
		for ; ; <-ticker.C {
			head, err := s.eth.BlockByNumber(ctx, nil)
			if err != nil {
				s.Service.Errors <- errors.Wrap(err, "failed to get block count")
				continue
			}

			s.Log.WithField("height", head.NumberU64()).Debug("fetched new head")

			for head.NumberU64()-s.config.Confirmations > cursor.Uint64() {
				s.blocksCh <- cursor.Uint64()
				cursor.Add(cursor, big.NewInt(1))
			}
		}
	}()
}

func (s *Service) processTXs(ctx context.Context) {
	for tx := range s.txCh {
		entry := s.Log.WithField("tx", tx.Hash().Hex())
		for {
			entry.Debug("processing tx")
			if err := s.processTX(ctx, tx); err != nil {
				s.Log.WithError(err).Error("failed to process tx")
				continue
			}
			entry.Debug("tx processed")
			break
		}

	}
}

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
	address := s.state.AddressAt(ctx, tx.Timestamp, tx.To().Hex())
	if address == nil {
		return nil
	}

	price := s.state.PriceAt(ctx, tx.Timestamp)
	if price == nil {
		s.Log.WithField("tx", tx.Hash().String()).Error("price is not set, skipping tx")
		return nil
	}

	receiver := s.state.Balance(ctx, *address)
	if receiver == nil {
		s.Log.WithField("tx", tx.Hash().String()).Error("balance not found, skipping tx")
		return nil
	}

	// amount = value * price / 10^18
	div := new(big.Int).Mul(big.NewInt(1000000000), big.NewInt(1000000000))
	bigPrice := big.NewInt(*price)

	amount := new(big.Int).Mul(tx.Value(), bigPrice)
	amount = amount.Div(amount, div)
	if !amount.IsUint64() {
		s.Log.WithField("tx", tx.Hash().String()).Error("amount overflow, skipping tx")
		return nil
	}

	s.Log.Info("submitting issuance request")

	reference := tx.Hash().Hex()
	// yoba eth hex trimming
	if len(reference) > 64 {
		reference = reference[len(reference)-64:]
	}
	request := s.CraftIssuanceRequest(supervisor.IssuanceRequestOpt{
		Reference: reference,
		Receiver:  *receiver,
		Amount:    amount.Uint64(),
		Details: resources.DepositDetails{
			Source: tx.Hash().Hex(),
			Price:  types.Amount(bigPrice.Int64()),
		}.Encode(),
	})

	if err = request.Submit(); err != nil {
		entry := s.Log.
			WithField("tx", tx.Hash().String()).
			WithField("block", tx.BlockNumber).
			WithError(err)

		if serr, ok := errors.Cause(err).(horizon.SubmitError); ok {
			opCodes := serr.OperationCodes()
			if len(opCodes) == 1 {
				switch opCodes[0] {
				// safe to move on
				case "op_reference_duplication":
					entry.Info("tx failed")
					return nil
				}
			}
		}

		entry.Error("failed to submit issuance request")
		return err
	}

	s.Log.Info("issuance request submitted")

	return nil
}
