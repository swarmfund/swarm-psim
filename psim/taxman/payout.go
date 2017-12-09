package taxman

import (
	"encoding/json"
	"fmt"
	"time"

	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/psim/psim/taxman/internal/payout"
	"gitlab.com/tokend/psim/psim/taxman/internal/resource"
	"gitlab.com/tokend/psim/psim/taxman/internal/snapshoter"
	"gitlab.com/tokend/psim/psim/utils"
	"gitlab.com/distributed_lab/sse-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func (s *Service) payoutStateLock(ledger int64) string {
	return fmt.Sprintf("%s/lock", s.payoutStateKey(ledger))
}

func (s *Service) payoutStateKey(ledger int64) string {
	return fmt.Sprintf("taxman/payouts/%d", ledger)
}

func (s *Service) GetPayoutState(ledger int64) (*resource.PayoutState, error) {
	kv, err := s.discovery.Get(s.payoutStateKey(ledger))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get kv")
	}

	if kv == nil || kv.Value == nil {
		return nil, nil
	}

	var state resource.PayoutState
	if err := json.Unmarshal(kv.Value, &state); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	return &state, nil
}

func (s *Service) SetPayoutState(ledger int64, state *resource.PayoutState) error {
	value, err := json.Marshal(&state)
	if err != nil {
		return errors.Wrap(err, "failed to marshal")
	}
	err = s.discovery.Set(discovery.KVPair{
		Key:   s.payoutStateKey(ledger),
		Value: value,
	})
	if err != nil {
		return errors.Wrap(err, "failed to set kv")
	}
	return nil
}

func (s *Service) processEvent(event sse.Event) (bool, error) {
	// parse horizon transaction response
	tx := horizon.Transaction{}
	err := json.NewDecoder(event.Data).Decode(&tx)
	if err != nil {
		return true, errors.Wrap(err, "failed to unmarshal event")
	}

	err = s.txHandler.Handle(tx)
	if err != nil {
		return true, errors.Wrap(err, "failed to process tx")
	}

	payoutPeriod := s.state.GetPayoutPeriod()
	if payoutPeriod == nil {
		return false, nil
	}

	isPayOutTime := tx.CreatedAt.After(s.NextPayout)

	if isPayOutTime {
		s.log.WithField("period", payoutPeriod.String()).Info("payout time")

		snapshot, err := snapshoter.NewBuilder(s.state).Build()
		if err != nil {
			// failure to build snapshoter results in unrecoverable state, aborting
			s.errors <- errors.Wrap(err, "failed to build snapshoter")
			return true, err
		}

		s.snapshots.Add(snapshot)

		// FIXME it's really bad
		// started as refactoring with some idea in mind then got distracted and
		// this happened, at that point it was too far to reset
		txBuilder := func() *horizon.TransactionBuilder {
			return s.horizon.Transaction(&horizon.TransactionBuilder{
				Source: s.config.Source,
			})
		}
		ask := func(request *resource.SyncRequest) (utils.AskResult, *utils.AskMeta) {
			return utils.AskNeighbors(s.discoveryService, request)
		}
		payout := payout.New(snapshot, payout.TransactionBuilder(txBuilder), s.config.Signer, payout.Ask(ask))

		for {
			if s.config.DisableVerify {
				s.log.WithField("ledger", snapshot.Ledger).Info("skipping verify")
				break
			}

			ctx, _ := context.WithTimeout(s.ctx, 5*time.Second)
			key := s.payoutStateKey(snapshot.Ledger)
			if err = payout.Ensure(s.discovery.Lock(ctx, key)); err != nil {
				s.errors <- errors.Wrap(err, "failed to ensure payout")
				continue
			}
			break
		}

		s.NextPayout = tx.CreatedAt.Add(*payoutPeriod)
		s.log.WithField("next_payout", s.NextPayout.String()).Info("successful payout")
	}

	s.TxCursor = tx.ID

	return false, nil
}
