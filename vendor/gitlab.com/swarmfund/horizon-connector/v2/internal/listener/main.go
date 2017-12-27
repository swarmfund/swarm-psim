package listener

import (
	"gitlab.com/swarmfund/horizon-connector/v2/internal/operation"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/resources"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/transaction"
)

type Q struct {
	tx *transaction.Q
	op *operation.Q
}

func NewQ(tx *transaction.Q, op *operation.Q) *Q {
	return &Q{
		tx, op,
	}
}

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

func (q *Q) Requests(result chan<- resources.Request) <-chan error {
	errs := make(chan error)
	go func() {
		defer func() {
			close(errs)
		}()
		cursor := ""
		for {
			requests, err := q.op.Requests(cursor)
			if err != nil {
				errs <- err
				continue
			}
			for _, request := range requests {
				result <- request
				cursor = request.PagingToken
			}
		}
	}()
	return errs
}
