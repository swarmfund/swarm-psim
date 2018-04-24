package lchanges

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
)

type TXStreamer interface {
	//StreamTransactions(ctx context.Context) (<-chan horizon.TransactionEvent, <-chan error)
	StreamTXs(ctx context.Context, stopOnEmptyPage bool) <-chan horizon.TXPacket
}

type TimedLedgerChange struct {
	Change xdr.LedgerEntryChange
	Time   time.Time
}

type Streamer struct {
	log             *logan.Entry
	txStreamer      TXStreamer
	stopOnEmptyPage bool

	timedChangesStream chan TimedLedgerChange
}

func NewStreamer(log *logan.Entry, txStreamer TXStreamer, stopOnEmptyPage bool) *Streamer {
	return &Streamer{
		log:             log.WithField("helper-runner", "ledger_changes_streamer"),
		txStreamer:      txStreamer,
		stopOnEmptyPage: stopOnEmptyPage,

		timedChangesStream: make(chan TimedLedgerChange),
	}
}

// GetStream returns stream of TimedLedgerChange where all the data is streamed.
// Consumers of Streamer should naturally listen to this channel as channels with work.
func (s Streamer) GetStream() <-chan TimedLedgerChange {
	return s.timedChangesStream
}

// Run is a blocking method, Run returns only if:
// - ctx cancelled;
// - txStream was closed (if stopOnEmptyPage is true)
//
// Run is not supposed to be called more than once - it closes LC stream in defer.
func (s *Streamer) Run(ctx context.Context) {
	s.log.Info("Started listening Transactions stream.")
	txStream := s.txStreamer.StreamTXs(ctx, s.stopOnEmptyPage)

	defer func() {
		close(s.timedChangesStream)
	}()

	// TODO Consider counting TXs per day and logging this number with each TX day log.
	var lastLoggedTXYearDay int
	for {
		select {
		case <-ctx.Done():
			return
		case txPacket, ok := <-txStream:
			if app.IsCanceled(ctx) {
				s.log.Info("Received cancel - stopping.")
				return
			}

			if !ok {
				// No more Transactions in the system
				s.log.Info("TX channel was closed - no more Transactions - stopping.")
				return
			}

			txEvent, err := txPacket.Unwrap()
			if err != nil {
				s.log.WithError(err).Error("TXStreamer sent error into its error channel.")
				continue
			}

			if txEvent.Transaction == nil {
				continue
			}

			if txEvent.Meta.LatestLedger.ClosedAt.YearDay() != lastLoggedTXYearDay {
				// New day TX
				s.log.WithField("tx_time", txEvent.Meta.LatestLedger.ClosedAt).Info("Received next day TX.")
				lastLoggedTXYearDay = txEvent.Meta.LatestLedger.ClosedAt.YearDay()
			}

			s.streamChanges(ctx, *txEvent.Transaction)
			continue
		}
	}
}

func (s *Streamer) streamChanges(ctx context.Context, tx horizon.Transaction) {
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
