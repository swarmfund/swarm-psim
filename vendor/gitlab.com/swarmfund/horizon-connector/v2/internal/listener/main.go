package listener

import (
	"gitlab.com/swarmfund/horizon-connector/v2/internal/operation"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/resources"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/transaction"
	"context"
)

type Q struct {
	txQ *transaction.Q
	// TODO Rename - it'a actually RequestQ
	opQ *operation.Q
}

func NewQ(tx *transaction.Q, op *operation.Q) *Q {
	return &Q{
		tx,
		op,
	}
}

// DEPRECATED use StreamAllReviewableRequests instead
func (q *Q) Requests(result chan<- resources.Request) <-chan error {
	errs := make(chan error)
	go func() {
		defer func() {
			close(errs)
		}()
		cursor := ""
		for {
			requests, err := q.opQ.AllRequests(cursor)
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
// DEPRECATED Use StreamWithdrawalRequests instead
func (q *Q) WithdrawalRequests(result chan<- resources.Request) <-chan error {
	errs := make(chan error)

	go func() {
		defer func() {
			close(errs)
		}()

		cursor := ""
		for {
			requests, err := q.opQ.WithdrawalRequests(cursor)
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

func (q *Q) StreamAllCheckSaleStateOps(ctx context.Context, buffer int) <-chan CheckSaleStateResponse {
	return streamCheckSaleState(q, ctx, buffer)
}

func (q *Q) StreamAllCreateKYCRequestOps(ctx context.Context, buffer int) <-chan CreateKYCRequestOpResponse {
	return streamCreateKYCRequestOp(q, ctx, buffer)
}
