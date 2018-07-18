package listener

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/regources"
)

// StreamTransactionsV2 streams transactions fetched for specified filters.
// If there is no new transactions, but ledger has been closed, `TransactionV2Event` with nil tx will be returned.
// Consumer should not rely on closing of any of this channels.
func (q *Q) StreamTransactionsV2(ctx context.Context, effects, entryTypes []int,
) (<-chan regources.TransactionV2Event, <-chan error) {
	txStream := make(chan regources.TransactionV2Event)
	errChan := make(chan error)

	go func() {
		cursor := ""
		for {
			select {
			case <-ctx.Done():
				return
			default:
				break
			}

			transactionsV2, meta, err := q.txV2Q.TransactionsByEffectsAndEntryTypes(cursor, effects, entryTypes)
			if err != nil {
				errChan <- errors.Wrap(err, "Failed to obtain Transactions", logan.F{
					"cursor" : cursor,
					"effects": effects,
					"entry_types" : entryTypes,
				})
				continue
			}

			for _, tx := range transactionsV2 {
				ohaigo := tx

				txEvent := regources.TransactionV2Event{
					TransactionV2: &ohaigo,
					// emulating discrete transactions stream by spoofing meta
					// to not let bump cursor too much before actually consuming all transactions
					Meta: regources.PageMeta{
						LatestLedger: regources.LedgerMeta{
							ClosedAt: tx.LedgerCloseTime,
						},
					},
				}
				ok := streamTxV2Event(ctx, txEvent, txStream)
				if !ok {
					// Ctx was canceled
					return
				}

				cursor = tx.PT
			}

			// letting consumer know about current ledger cursor
			ok := streamTxV2Event(ctx, regources.TransactionV2Event{
				TransactionV2: nil,
				Meta:          *meta,
			}, txStream)
			if !ok {
				// Ctx was canceled
				return
			}
		}
	}()

	return txStream, errChan
}

func streamTxV2Event(ctx context.Context, txEvent regources.TransactionV2Event,
txStream chan<- regources.TransactionV2Event) bool {
	select {
	case <-ctx.Done():
		return false
	case txStream <- txEvent:
		return true
	}
}