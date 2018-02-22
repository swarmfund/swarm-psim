package listener


import (
	"context"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
)

// DEPRECATED Use StreamTransactions instead
func (q *Q) Transactions(result chan<- resources.TransactionEvent) <-chan error {
	errs := make(chan error)
	go func() {
		defer func() {
			close(errs)
		}()
		cursor := ""
		for {
			transactions, meta, err := q.tx.Transactions(cursor)
			if err != nil {
				errs <- err
				continue
			}
			for _, tx := range transactions {
				ohaigo := tx
				result <- resources.TransactionEvent{
					Transaction: &ohaigo,
					// emulating discrete transactions stream by spoofing meta
					// to not let bump cursor too much before actually consuming all transactions
					Meta: resources.PageMeta{
						LatestLedger: resources.LedgerMeta{
							ClosedAt: tx.CreatedAt,
						},
					},
				}
				cursor = tx.PagingToken
			}
			// letting consumer know about current ledger cursor
			result <- resources.TransactionEvent{
				Transaction: nil,
				Meta:        *meta,
			}
		}
	}()
	return errs
}

func (q *Q) StreamTransactions(ctx context.Context) (<-chan resources.TransactionEvent, <- chan error) {
	txStream := make(chan resources.TransactionEvent)
	errChan := make(chan error)

	go func() {
		defer func() {
			close(txStream)
			close(errChan)
		}()

		cursor := ""
		for {
			transactions, meta, err := q.tx.Transactions(cursor)
			if err != nil {
				errChan <- errors.Wrap(err, "Failed to obtain Transactions", logan.F{"cursor": cursor})
				continue
			}

			for _, tx := range transactions {
				ohaigo := tx

				txStream <- resources.TransactionEvent{
					Transaction: &ohaigo,
					// emulating discrete transactions stream by spoofing meta
					// to not let bump cursor too much before actually consuming all transactions
					Meta: resources.PageMeta{
						LatestLedger: resources.LedgerMeta{
							ClosedAt: tx.CreatedAt,
						},
					},
				}

				cursor = tx.PagingToken
			}

			// letting consumer know about current ledger cursor
			txStream <- resources.TransactionEvent{
				Transaction: nil,
				Meta:        *meta,
			}
		}
	}()

	return txStream, errChan
}
