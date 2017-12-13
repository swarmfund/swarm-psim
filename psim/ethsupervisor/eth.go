package ethsupervisor

import (
	"math/big"
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/psim/psim/ethsupervisor/internal"
	"fmt"
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
				BlockNumber: block.NumberU64(),
				Transaction: *tx,
			}
		}
	}
}

func (s *Service) watchHeight() {
	// TODO config
	cursor := *big.NewInt(2258360)
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for ;; <-ticker.C{
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
				//s.Log.WithField("cursor", cursor.Uint64()).Debug("cursor bumped")
			}
		}
	}()
}

func (s *Service) processTXs() {
	for tx := range s.txCh {
		s.processTX(tx)
	}
}

func (s *Service) processTX(tx internal.Transaction) {
	if tx.Value().Cmp(&s.depositThreshold) == -1 {
		return
	}

	if tx.To() == nil {
		return
	}

	address := s.state.AddressAt(tx.Timestamp, tx.To().String())
	if address == nil {
		return
	}

	// TODO get balance by address
	balanceID := "BBS5KRCNZZR2MRKXMJU2SAAXYFBTSJARCKNZVKTDWSIA62SQIERJP2GX"

	// TODO convert based on ONE and asset pair rate
	// 10^18/10^6
	var div int64 = 1000000000000
	amount := new(big.Int).Div(tx.Value(), big.NewInt(div)).Uint64()

	s.Log.Info("submitting issuance request")

	err := s.horizon.Transaction(&horizon.TransactionBuilder{Source: s.config.Supervisor.ExchangeKP}).
		Op(&horizon.CreateIssuanceRequestOp{
			Reference: tx.Hash().String(),
			Amount:    amount,
			Asset: "SUN",
			Receiver: balanceID,
		}).
		Sign(s.config.Supervisor.SignerKP).
		Submit()
	if err != nil {
		if serr, ok := errors.Cause(err).(horizon.SubmitError); ok {
			fmt.Println(string(serr.ResponseBody()))
		}
		s.Log.
			WithField("tx", tx.Hash().String()).
			WithField("block", tx.BlockNumber).
			WithError(err).
			Error("failed to submit issuance request")
	}
}
