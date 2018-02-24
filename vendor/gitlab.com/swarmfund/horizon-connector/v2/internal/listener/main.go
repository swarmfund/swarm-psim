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
		tx,
		op,
	}
}

// DEPRECATED Does not work any more. Can now only stream WithdrawalRequests.
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

// TODO Consider working with *Withdrawal* specific types.
func (q *Q) WithdrawalRequests(result chan<- resources.Request) <-chan error {
	errs := make(chan error)

	go func() {
		defer func() {
			close(errs)
		}()

		cursor := ""
		for {
			requests, err := q.op.WithdrawalRequests(cursor)
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
