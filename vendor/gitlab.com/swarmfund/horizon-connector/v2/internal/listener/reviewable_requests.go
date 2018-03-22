package listener

import (
	"context"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/resources"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/operation"
)

func (q *Q) StreamAllReviewableRequests(ctx context.Context) (<-chan ReviewableRequestEvent) {
	reqGetter := func(cursor string) ([]resources.Request, error) {
		return q.opQ.AllRequests(cursor)
	}

	return streamReviewableRequests(ctx, reqGetter)
}

// TODO When KYC is ready
//func (q *Q) StreamKYCRequests(ctx context.Context) (<-chan ReviewableRequestEvent) {
//	reqGetter := func(cursor string) ([]resources.Request, error) {
//		return q.opQ.Requests(cursor, operation.KYCReviewableRequestType)
//	}
//
//	return streamReviewableRequests(ctx, reqGetter)
//}

func (q *Q) StreamWithdrawalRequests(ctx context.Context) (<-chan ReviewableRequestEvent) {
	reqGetter := func(cursor string) ([]resources.Request, error) {
		return q.opQ.Requests(cursor, operation.WithdrawalsReviewableRequestType)
	}

	return streamReviewableRequests(ctx, reqGetter)
}

func streamReviewableRequests(ctx context.Context, reqGetter func(cursor string) ([]resources.Request, error)) (<-chan ReviewableRequestEvent) {
	reqStream := make(chan ReviewableRequestEvent)

	go func() {
		defer func() {
			close(reqStream)
		}()

		cursor := ""
		for {
			select {
			case <-ctx.Done():
				return
			default:
				break
			}

			//requests, err := q.opQ.AllRequests(cursor)
			requests, err := reqGetter(cursor)
			if err != nil {
				streamRequestEvent(ctx, ReviewableRequestEvent{
					body: nil,
					err:  err,
				}, reqStream)
				continue
			}

			for _, req := range requests {
				ohaigo := req

				reqEvent := ReviewableRequestEvent{
					body: &ohaigo,
					err:  err,
				}
				ok := streamRequestEvent(ctx, reqEvent, reqStream)
				if !ok {
					// Ctx was canceled
					return
				}

				cursor = req.PagingToken
			}
		}
	}()

	return reqStream
}

func streamRequestEvent(ctx context.Context, reqEvent ReviewableRequestEvent, reqStream chan<- ReviewableRequestEvent) bool {
	select {
	case <- ctx.Done():
		return false
	case reqStream <- reqEvent:
		return true
	}
}
