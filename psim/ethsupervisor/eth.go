package ethsupervisor

import (
	"math/big"

	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/ethsupervisor/internal"
)

// TODO defer
func (s *Service) processBlocks() {
	for blockNumber := range s.blocksCh {
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
				Transaction: *tx,
			}
		}
	}
}

func (s *Service) watchHeight() {
	cursor := *big.NewInt(0)
	go func() {
		for {
			head, err := s.eth.BlockByNumber(s.Ctx, nil)
			if err != nil {
				s.Service.Errors <- errors.Wrap(err, "failed to get block count")
				continue
			}

			s.Log.WithField("height", head.NumberU64()).Debug("fetched new head")

			// FIXME Magic 12 number
			for head.NumberU64()-12 > cursor.Uint64() {
				s.blocksCh <- cursor.Uint64()
				cursor.Add(&cursor, big.NewInt(1))
				s.Log.WithField("cursor", cursor.Uint64()).Debug("cursor bumped")
			}
		}
	}()
}

func (s *Service) processTXs() {
	for tx := range s.txCh {
		if tx.Value().Cmp(&s.depositTreshold) == -1 {
			continue
		}

		if tx.To() == nil {
			continue
		}

		address := s.state.AddressAt(tx.Timestamp, tx.To().String())
		if address != nil {
			continue
		}

		// TODO craft CERs
	}
}
