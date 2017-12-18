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

func (q *Q) Transactions(result chan<- resources.Transaction) <-chan error {
	errs := make(chan error)
	go func() {
		defer func() {
			close(errs)
		}()
		cursor := ""
		for {
			transactions, err := q.tx.Transactions(cursor)
			if err != nil {
				errs <- err
				continue
			}
			for _, tx := range transactions {
				result <- tx
				cursor = tx.PagingToken
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
