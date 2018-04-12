package notifier

import (
	"context"
	"fmt"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"strconv"
)

type ReviewableRequestConnector interface {
	GetRequestByID(requestID uint64) (*horizon.Request, error)
}

type ReviewedKYCRequestNotifier struct {
	approvedKYCEmailSender EmailSender
	rejectedKYCEmailSender EmailSender
	approvedRequestConfig  EventConfig
	rejectedRequestConfig  EventConfig
	requestConnector       ReviewableRequestConnector
	userConnector          UserConnector

	reviewRequestOpResponses <-chan horizon.ReviewRequestOpResponse
}

func (n *ReviewedKYCRequestNotifier) listenAndProcessReviewedKYCRequests(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case reviewRequestOpResponse, ok := <-n.reviewRequestOpResponses:
			if !ok {
				return nil
			}

			reviewRequestOp, err := reviewRequestOpResponse.Unwrap()
			if err != nil {
				return errors.Wrap(err, "ReviewRequestOpListener sent error")
			}

			cursor, err := strconv.ParseUint(reviewRequestOp.PT, 10, 64)
			if err != nil {
				return errors.Wrap(err, "failed to parse paging token", logan.F{
					"paging_token": reviewRequestOp.PT,
				})
			}

			if n.canNotifyAboutApprovedKYC(cursor) && n.isFullyApprovedKYC(*reviewRequestOp) {
				err := n.notifyAboutApprovedKYCRequest(ctx, reviewRequestOp.RequestID)
				if err != nil {
					return errors.Wrap(err, "failed to notify about approved KYC request", logan.F{
						"request_id": reviewRequestOp.RequestID,
					})
				}
			}

			if n.canNotifyAboutRejectedKYC(cursor) && n.isRejectedKYC(*reviewRequestOp) {
				err := n.notifyAboutRejectedKYCRequest(ctx, reviewRequestOp.RequestID, reviewRequestOp.ID)
				if err != nil {
					return errors.Wrap(err, "failed to notify about rejected KYC request", logan.F{
						"request_id": reviewRequestOp.RequestID,
					})
				}
			}
		}
	}
}

func (n *ReviewedKYCRequestNotifier) notifyAboutApprovedKYCRequest(ctx context.Context, requestID uint64) error {
	request, err := n.requestConnector.GetRequestByID(requestID)
	if err != nil {
		return errors.Wrap(err, "failed to get reviewable request", logan.F{
			"request_id": requestID,
		})
	}

	user, err := n.userConnector.User(request.Details.KYC.AccountToUpdateKYC)
	if err != nil {
		return errors.Wrap(err, "failed to load user", logan.F{
			"account_id": request.Details.KYC.AccountToUpdateKYC,
		})
	}
	if user == nil {
		return nil
	}

	emailAddress := user.Attributes.Email
	emailUniqueToken := n.buildApprovedKYCUniqueToken(emailAddress, request.Details.KYC.AccountToUpdateKYC, requestID)

	data := struct {
		Link string
	}{
		Link: n.approvedRequestConfig.Emails.TemplateLinkURL,
	}

	err = n.approvedKYCEmailSender.SendEmail(ctx, emailAddress, emailUniqueToken, data)
	if err != nil {
		return errors.Wrap(err, "failed to send email")
	}

	return nil
}

func (n *ReviewedKYCRequestNotifier) notifyAboutRejectedKYCRequest(ctx context.Context, requestID uint64, operationID string) error {
	request, err := n.requestConnector.GetRequestByID(requestID)
	if err != nil {
		return errors.Wrap(err, "failed to get reviewable request", logan.F{
			"request_id": requestID,
		})
	}

	user, err := n.userConnector.User(request.Details.KYC.AccountToUpdateKYC)
	if err != nil {
		return errors.Wrap(err, "failed to load user", logan.F{
			"account_id": request.Details.KYC.AccountToUpdateKYC,
		})
	}
	if user == nil {
		return nil
	}

	emailAddress := user.Attributes.Email
	emailUniqueToken := n.buildRejectedKYCUniqueToken(emailAddress, request.Details.KYC.AccountToUpdateKYC, request.Details.KYC.SequenceNumber, requestID)

	data := struct {
		Link string
	}{
		Link: n.rejectedRequestConfig.Emails.TemplateLinkURL,
	}

	err = n.rejectedKYCEmailSender.SendEmail(ctx, emailAddress, emailUniqueToken, data)
	if err != nil {
		return errors.Wrap(err, "failed to send email")
	}

	return nil
}

func (n *ReviewedKYCRequestNotifier) isFullyApprovedKYC(reviewRequestOp horizon.ReviewRequestOp) bool {
	if reviewRequestOp.Action != xdr.ReviewRequestOpActionApprove.ShortString() {
		return false
	}

	if reviewRequestOp.RequestType != xdr.ReviewableRequestTypeUpdateKyc.ShortString() {
		return false
	}

	return reviewRequestOp.IsFulfilled
}

func (n *ReviewedKYCRequestNotifier) isRejectedKYC(reviewRequestOp horizon.ReviewRequestOp) bool {
	if reviewRequestOp.Action != xdr.ReviewRequestOpActionReject.ShortString() {
		return false
	}

	if reviewRequestOp.RequestType != xdr.ReviewableRequestTypeUpdateKyc.ShortString() {
		return false
	}

	return true
}

func (n *ReviewedKYCRequestNotifier) canNotifyAboutApprovedKYC(cursor uint64) bool {
	return cursor >= n.approvedRequestConfig.Cursor
}

func (n *ReviewedKYCRequestNotifier) canNotifyAboutRejectedKYC(cursor uint64) bool {
	return cursor >= n.rejectedRequestConfig.Cursor
}

func (n *ReviewedKYCRequestNotifier) buildApprovedKYCUniqueToken(emailAddress, accountToUpdateKYC string, requestID uint64) string {
	return fmt.Sprintf("%s:%s:%d:%s", emailAddress, accountToUpdateKYC, requestID, n.approvedRequestConfig.Emails.RequestTokenSuffix)
}

func (n *ReviewedKYCRequestNotifier) buildRejectedKYCUniqueToken(emailAddress, accountToUpdateKYC string, kycSequence uint32, requestID uint64) string {
	return fmt.Sprintf("%s:%s:%d:%d:%s", emailAddress, accountToUpdateKYC, kycSequence, requestID, n.rejectedRequestConfig.Emails.RequestTokenSuffix)
}
