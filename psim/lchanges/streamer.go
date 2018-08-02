package lchanges

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/internal"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/regources"
)

type TXStreamer interface {
	StreamTXsFromCursor(ctx context.Context, cursor string, stopOnEmptyPage bool) <-chan horizon.TXPacket
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
//
// Cursor is PagingToken for Transaction to start from,
// use empty string to start from the very beginning of Transactions history.
func (s *Streamer) Run(ctx context.Context, cursor string) {
	s.log.Info("Started listening Transactions stream.")
	txStream := s.txStreamer.StreamTXsFromCursor(ctx, "", s.stopOnEmptyPage)

	defer func() {
		close(s.timedChangesStream)
	}()

	var lastLoggedTXYearDay int
	var txsPerDay uint64
	for {
		select {
		case <-ctx.Done():
			return
		case txPacket, ok := <-txStream:
			if running.IsCancelled(ctx) {
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

			txsPerDay += 1

			if txEvent.Meta.LatestLedger.ClosedAt.YearDay() != lastLoggedTXYearDay {
				// New day TX
				s.log.WithFields(logan.F{
					"tx_time":              txEvent.Meta.LatestLedger.ClosedAt,
					"txs_per_previous_day": txsPerDay,
				}).Info("Received next day TX.")

				lastLoggedTXYearDay = txEvent.Meta.LatestLedger.ClosedAt.YearDay()
				txsPerDay = 0
			}

			s.streamChanges(ctx, *txEvent.Transaction)
			continue
		}
	}
}

func (s *Streamer) streamChanges(ctx context.Context, tx regources.Transaction) {
	for _, change := range internal.LedgerChanges(&tx) {
		timedChange := TimedLedgerChange{
			Change: change,
			Time:   tx.LedgerCloseTime,
		}

		select {
		case <-ctx.Done():
			return
		case s.timedChangesStream <- timedChange:
			continue
		}
	}
}
