package ethsupervisor

import (
	"math/big"
	"time"

	"strings"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/psim/psim/ethsupervisor/internal"
)

// TODO defer
func (s *Service) processBlocks() {
	for blockNumber := range s.blocksCh {
		s.Log.WithField("number", blockNumber).Debug("processing block")
		block, err := s.eth.BlockByNumber(s.Ctx, big.NewInt(int64(blockNumber)))
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

func (s *Service) watchHeight() {
	// TODO config
	cursor := *big.NewInt(2271294)
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for ; ; <-ticker.C {
			head, err := s.eth.BlockByNumber(s.Ctx, nil)
			if err != nil {
				s.Service.Errors <- errors.Wrap(err, "failed to get block count")
				continue
			}

			s.Log.WithField("height", head.NumberU64()).Debug("fetched new head")

			// FIXME Magic number
			for head.NumberU64()-2 > cursor.Uint64() {
				s.blocksCh <- cursor.Uint64()
				cursor.Add(&cursor, big.NewInt(1))
				//s.Log.WithField("cursor", cursor.Uint64()).Debug("cursor bumped")
			}
		}
	}()
}

func (s *Service) processTXs() {
	for tx := range s.txCh {
		for {
			if err := s.processTX(tx); err != nil {
				s.Log.WithError(err).Error("failed to process tx")
				continue
			}
			break
		}

	}
}

func (s *Service) processTX(tx internal.Transaction) (err error) {
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
	address := s.state.AddressAt(tx.Timestamp, strings.ToLower(tx.To().String()))
	if address == nil {
		return nil
	}

	price := s.state.PriceAt(tx.Timestamp)
	if price == nil {
		s.Log.WithField("tx", tx.Hash().String()).Error("price is not set, skipping tx")
		return nil
	}

	receiver := s.state.Balance(*address)
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

	if err = s.PrepareCERTx(reference, *receiver, amount.Uint64()).Submit(); err != nil {
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
