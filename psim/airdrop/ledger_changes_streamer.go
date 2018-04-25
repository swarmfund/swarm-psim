package airdrop

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/swarmfund/psim/psim/app"
)

type TXStreamer interface {
	StreamTransactions(ctx context.Context) (<-chan horizon.TransactionEvent, <-chan error)
}

type TimedLedgerChange struct {
	Change xdr.LedgerEntryChange
	Time   time.Time
}

// DEPRECATED
// Use lchanges.Streamer instead
type LedgerChangesStreamer struct {
	log        *logan.Entry
	txStreamer TXStreamer

	timedChangesStream chan TimedLedgerChange
}

// DEPRECATED
// Use lchanges.NewStreamer instead
func NewLedgerChangesStreamer(log *logan.Entry, txStreamer TXStreamer) *LedgerChangesStreamer {
	return &LedgerChangesStreamer{
		log:        log.WithField("helper-runner", "ledger_changes_streamer"),
		txStreamer: txStreamer,

		timedChangesStream: make(chan TimedLedgerChange),
	}
}

// DEPRECATED
// Use lchanges.Streamer instead
func (s *LedgerChangesStreamer) Run(ctx context.Context) <-chan TimedLedgerChange {
	s.log.Info("Started listening Transactions stream.")
	txStream, txStreamerErrs := s.txStreamer.StreamTransactions(ctx)

	var isFirstTX = true
	var lastLoggedTXYearDay int
	go app.RunOverIncrementalTimer(ctx, s.log, "ledger_changes_processor", func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return nil
		case txEvent := <-txStream:
			if app.IsCanceled(ctx) {
				return nil
			}

			if txEvent.Transaction == nil {
				return nil
			}

			if isFirstTX {
				s.log.WithField("tx_time", txEvent.Meta.LatestLedger.ClosedAt).Info("Received first TX.")
				lastLoggedTXYearDay = txEvent.Meta.LatestLedger.ClosedAt.YearDay()
				isFirstTX = false
			} else {
				if txEvent.Meta.LatestLedger.ClosedAt.YearDay() != lastLoggedTXYearDay {
					// New day TX
					s.log.WithField("tx_time", txEvent.Meta.LatestLedger.ClosedAt).Info("Received next day TX.")
					lastLoggedTXYearDay = txEvent.Meta.LatestLedger.ClosedAt.YearDay()
				}
			}

			s.streamChanges(ctx, *txEvent.Transaction)
			return nil
		case txStreamerErr := <-txStreamerErrs:
			s.log.WithError(txStreamerErr).Error("TXStreamer sent error into its error channel.")
			return nil
		}
	}, 0, 10*time.Second)

	return s.timedChangesStream
}

func (s *LedgerChangesStreamer) streamChanges(ctx context.Context, tx horizon.Transaction) {
	for _, change := range tx.LedgerChanges() {
		timedChange := TimedLedgerChange{
			Change: change,
			Time:   tx.CreatedAt,
		}

		select {
		case <-ctx.Done():
			return
		case s.timedChangesStream <- timedChange:
			continue
		}
	}
}
