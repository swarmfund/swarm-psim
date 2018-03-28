package listener

import (
	"context"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/resources"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/operation"
	"time"
	"fmt"
)

func (q *Q) StreamAllReviewableRequests(ctx context.Context) (<-chan ReviewableRequestEvent) {
	reqGetter := func(cursor string) ([]resources.Request, error) {
		return q.opQ.AllRequests(cursor)
	}

	return streamReviewableRequests(ctx, reqGetter, false)
}

// StreamAllKYCRequests streams all ReviewableRequests of type KYC from very beginning,
// sorted by ID
func (q *Q) StreamAllKYCRequests(ctx context.Context, endlessly bool) (<-chan ReviewableRequestEvent) {
	return q.streamKYCRequests(ctx, "", !endlessly)
}

func (q *Q) StreamKYCRequestsUpdatedAfter(ctx context.Context, updatedAfter time.Time, endlessly bool) (<-chan ReviewableRequestEvent) {
	filters := fmt.Sprintf("updated_after=%d", updatedAfter.UTC().Unix())
	return q.streamKYCRequests(ctx, filters, !endlessly)
}

func (q *Q) streamKYCRequests(ctx context.Context, filters string, stopOnEmptyPage bool) (<-chan ReviewableRequestEvent) {
	reqGetter := func(cursor string) ([]resources.Request, error) {
		return q.opQ.Requests(filters, cursor, operation.KYCReviewableRequestType)
	}

	return streamReviewableRequests(ctx, reqGetter, stopOnEmptyPage)
}

func (q *Q) StreamWithdrawalRequests(ctx context.Context) (<-chan ReviewableRequestEvent) {
	reqGetter := func(cursor string) ([]resources.Request, error) {
		return q.opQ.Requests("", cursor, operation.WithdrawalsReviewableRequestType)
	}

	return streamReviewableRequests(ctx, reqGetter, false)
}

func streamReviewableRequests(ctx context.Context, reqGetter func(cursor string) ([]resources.Request, error), stopOnEmptyPage bool) (<-chan ReviewableRequestEvent) {
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

			reqBatch, err := reqGetter(cursor)
			if err != nil {
				streamRequestEvent(ctx, ReviewableRequestEvent{
					body: nil,
					err:  err,
				}, reqStream)
				continue
			}

			if stopOnEmptyPage || len(reqBatch) == 0 {
				// The stream channel is closed in defer.
				return
			}

			for _, req := range reqBatch {
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
