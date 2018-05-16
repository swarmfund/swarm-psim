package listener

import (
	"context"
	"gitlab.com/tokend/horizon-connector/internal/resources"
	"gitlab.com/tokend/horizon-connector/internal/operation"
	"time"
	"fmt"
)

func (q *Q) StreamAllReviewableRequests(ctx context.Context) (<-chan ReviewableRequestEvent) {
	reqGetter := func(cursor string) ([]resources.Request, error) {
		return q.opQ.AllRequests(cursor)
	}

	return streamReviewableRequests(ctx, reqGetter, "", false)
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
	return q.getAndStreamReviewableRequests(ctx, filters, "", operation.KYCReviewableRequestType, stopOnEmptyPage)
}

// StreamWithdrawalRequests streams all ReviewableRequests of type Withdraw and TwoStepWithdraw
func (q *Q) StreamWithdrawalRequests(ctx context.Context) (<-chan ReviewableRequestEvent) {
	return q.getAndStreamReviewableRequests(ctx, "", "", operation.WithdrawalsReviewableRequestType, false)
}

// StreamWithdrawalRequestsOfAsset streams all Withdraw and TwoStepWithdraw ReviewableRequests
// with filter by provided destAssetCode
func (q *Q) StreamWithdrawalRequestsOfAsset(ctx context.Context, destAssetCode string, reverseOrder bool) (<-chan ReviewableRequestEvent) {
	getParams := fmt.Sprintf("dest_asset_code=%s", destAssetCode)

	if reverseOrder {
		getParams += "&order=desc"
	}

	return q.getAndStreamReviewableRequests(ctx, getParams, "", operation.WithdrawalsReviewableRequestType, false)
}

func (q *Q) getAndStreamReviewableRequests(ctx context.Context, getParams, cursor string, reqType operation.ReviewableRequestType, stopOnEmptyPage bool) (<-chan ReviewableRequestEvent) {
	reqGetter := func(cursor string) ([]resources.Request, error) {
		return q.opQ.Requests(getParams, cursor, reqType)
	}

	return streamReviewableRequests(ctx, reqGetter, cursor, stopOnEmptyPage)
}

func streamReviewableRequests(ctx context.Context, reqGetter func(cursor string) ([]resources.Request, error), cursor string, stopOnEmptyPage bool) (<-chan ReviewableRequestEvent) {
	reqStream := make(chan ReviewableRequestEvent)

	go func() {
		defer func() {
			close(reqStream)
		}()

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

			if stopOnEmptyPage && len(reqBatch) == 0 {
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
