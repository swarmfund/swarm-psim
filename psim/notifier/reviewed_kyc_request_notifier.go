package notifier

import (
	"context"
	"fmt"
	"strconv"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/kyc"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/regources"
)

type ReviewableRequestConnector interface {
	GetRequestByID(requestID uint64) (*regources.ReviewableRequest, error)
}

type ReviewedKYCRequestNotifier struct {
	log                    *logan.Entry
	approvedKYCEmailSender EmailSender
	usaKYCEmailSender      EmailSender
	rejectedKYCEmailSender EmailSender
	approvedRequestConfig  EventConfig
	usaKYCConfig           EventConfig
	rejectedRequestConfig  EventConfig
	requestConnector       ReviewableRequestConnector
	userConnector          UserConnector
	kycDataHelper          KYCDataHelper

	reviewRequestOpResponses <-chan horizon.ReviewRequestOpResponse
}

func (n *ReviewedKYCRequestNotifier) listenAndProcessReviewedKYCRequests(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	case reviewRequestOpResponse, ok := <-n.reviewRequestOpResponses:
		if !ok {
			return nil
		}

		reviewRequestOp, err := reviewRequestOpResponse.Unwrap()
		if err != nil {
			return errors.Wrap(err, "ReviewRequestOpStreamer sent error")
		}

		fields := logan.F{
			"request_id":   reviewRequestOp.RequestID,
			"paging_token": reviewRequestOp.PT,
		}

		n.log.WithFields(fields).Debug("processing request")

		if reviewRequestOp.RequestType != xdr.ReviewableRequestTypeUpdateKyc.ShortString() {
			// Normally should never happen, but just in case. Ignoring.
			return nil
		}

		cursor, err := strconv.ParseUint(reviewRequestOp.PT, 10, 64)
		if err != nil {
			return errors.Wrap(err, "Failed to parse PagingToken", fields)
		}

		isCanNotifyApproved := n.canNotifyAboutApprovedKYC(cursor)
		isCanNotifyRejected := n.canNotifyAboutRejectedKYC(cursor)
		isFullyApproved := n.isFullyApprovedKYC(*reviewRequestOp)
		isRejected := n.isRejectedKYC(*reviewRequestOp)

		fields.Merge(logan.F{
			"is_can_notify_approved": isCanNotifyApproved,
			"is_can_notify_rejected": isCanNotifyRejected,
			"is_fully_approved":      isFullyApproved,
			"is_rejected":            isRejected,
		})

		if isCanNotifyApproved && isFullyApproved {
			n.log.WithFields(fields).Debug("notifying approve")
			err := n.notifyAboutApprovedKYCRequest(ctx, reviewRequestOp.RequestID)
			if err != nil {
				return errors.Wrap(err, "failed to notify about approved KYC request", fields)
			}
		}

		if isCanNotifyRejected && isRejected {
			n.log.WithFields(fields).Debug("notifying reject")
			err := n.notifyAboutRejectedKYCRequest(ctx, reviewRequestOp.RequestID)
			if err != nil {
				return errors.Wrap(err, "Failed to notify about rejected KYCRequest", fields)
			}
		}

		return nil
	}
}

func (n *ReviewedKYCRequestNotifier) notifyAboutApprovedKYCRequest(ctx context.Context, requestID uint64) error {
	request, err := n.requestConnector.GetRequestByID(requestID)
	if err != nil {
		return errors.Wrap(err, "failed to get reviewable request", logan.F{
			"request_id": requestID,
		})
	}

	kycRequest := request.Details.KYC

	if kycRequest.AccountTypeToSet.Int != int(xdr.AccountTypeGeneral) {
		return nil
	}

	user, err := n.userConnector.User(kycRequest.AccountToUpdateKYC)
	if err != nil {
		return errors.Wrap(err, "failed to load user", logan.F{
			"account_id": kycRequest.AccountToUpdateKYC,
		})
	}
	if user == nil {
		return nil
	}

	emailAddress := user.Attributes.Email
	emailUniqueToken := n.buildApprovedKYCUniqueToken(emailAddress, kycRequest.AccountToUpdateKYC, requestID)

	kycFirstName, err := n.kycDataHelper.getKYCFirstName(kycRequest.KYCData)
	if err != nil {
		return errors.Wrap(err, "failed to get blob KYC data")
	}

	data := struct {
		Link      string
		FirstName string
	}{
		Link:      n.approvedRequestConfig.Emails.TemplateLinkURL,
		FirstName: kycFirstName,
	}

	err = n.approvedKYCEmailSender.SendEmail(ctx, emailAddress, emailUniqueToken, data)
	if err != nil {
		return errors.Wrap(err, "failed to send email")
	}

	return nil
}

func (n *ReviewedKYCRequestNotifier) notifyAboutRejectedKYCRequest(ctx context.Context, requestID uint64) error {
	request, err := n.requestConnector.GetRequestByID(requestID)
	if err != nil {
		return errors.Wrap(err, "failed to get reviewable request", logan.F{
			"request_id": requestID,
		})
	}

	kycRequest := request.Details.KYC

	if kycRequest.AccountTypeToSet.Int != int(xdr.AccountTypeGeneral) {
		return nil
	}

	user, err := n.userConnector.User(kycRequest.AccountToUpdateKYC)
	if err != nil {
		return errors.Wrap(err, "failed to load user", logan.F{
			"account_id": kycRequest.AccountToUpdateKYC,
		})
	}
	if user == nil {
		return nil
	}

	kycFirstName, err := n.kycDataHelper.getKYCFirstName(kycRequest.KYCData)
	if err != nil {
		return errors.Wrap(err, "Failed to get Blob KYCData")
	}

	emailAddress := user.Attributes.Email
	emailUniqueToken := n.buildRejectedKYCUniqueToken(emailAddress, kycRequest.AccountToUpdateKYC, kycRequest.SequenceNumber, requestID)

	data := struct {
		Link         string
		FirstName    string
		RejectReason string
	}{
		Link:         n.rejectedRequestConfig.Emails.TemplateLinkURL,
		FirstName:    kycFirstName,
		RejectReason: request.RejectReason,
	}

	err = n.rejectedKYCEmailSender.SendEmail(ctx, emailAddress, emailUniqueToken, data)
	if err != nil {
		return errors.Wrap(err, "failed to send email")
	}

	return nil
}

func (n *ReviewedKYCRequestNotifier) tryNotifyAboutUSAKyc(ctx context.Context, requestID uint64) error {
	request, err := n.requestConnector.GetRequestByID(requestID)
	if err != nil {
		return errors.Wrap(err, "Failed to get ReviewableRequest from Horizon", logan.F{
			"request_id": requestID,
		})
	}
	fields := logan.F{
		"request": request,
	}

	kycRequest := request.Details.KYC

	if kycRequest.AllTasks&kyc.TaskUSA == 0 {
		// Not a USA User, at least it's not marked as USA user yet.
		return nil
	}

	user, err := n.userConnector.User(kycRequest.AccountToUpdateKYC)
	if err != nil {
		return errors.Wrap(err, "Failed to load User from Horizon", fields)
	}
	if user == nil {
		return errors.From(errors.New("No User found in Horizon"), fields)
	}
	fields["user"] = user

	emailAddr := user.Attributes.Email
	emailUniqueToken := emailAddr + n.usaKYCConfig.Emails.RequestTokenSuffix

	kycFirstName, err := n.kycDataHelper.getKYCFirstName(kycRequest.KYCData)
	if err != nil {
		return errors.Wrap(err, "Failed to obtain Blob KYCData")
	}

	templateData := struct {
		Link      string
		FirstName string
	}{
		Link:      n.approvedRequestConfig.Emails.TemplateLinkURL,
		FirstName: kycFirstName,
	}

	err = n.usaKYCEmailSender.SendEmail(ctx, emailAddr, emailUniqueToken, templateData)
	if err != nil {
		return errors.Wrap(err, "Failed to send email")
	}

	return nil
}

func (n *ReviewedKYCRequestNotifier) isFullyApprovedKYC(reviewRequestOp horizon.ReviewRequestOp) bool {
	if reviewRequestOp.Action != xdr.ReviewRequestOpActionApprove.ShortString() {
		return false
	}

	return reviewRequestOp.IsFulfilled
}

func (n *ReviewedKYCRequestNotifier) isRejectedKYC(reviewRequestOp horizon.ReviewRequestOp) bool {
	if reviewRequestOp.Action != xdr.ReviewRequestOpActionReject.ShortString() {
		return false
	}

	return true
}

func (n *ReviewedKYCRequestNotifier) canNotifyAboutApprovedKYC(cursor uint64) bool {
	return cursor >= n.approvedRequestConfig.Cursor && !n.approvedRequestConfig.Disabled
}

func (n *ReviewedKYCRequestNotifier) canNotifyAboutUSAKyc(cursor uint64) bool {
	return cursor >= n.usaKYCConfig.Cursor
}

func (n *ReviewedKYCRequestNotifier) canNotifyAboutRejectedKYC(cursor uint64) bool {
	return cursor >= n.rejectedRequestConfig.Cursor && !n.rejectedRequestConfig.Disabled
}

func (n *ReviewedKYCRequestNotifier) buildApprovedKYCUniqueToken(emailAddress, accountToUpdateKYC string, requestID uint64) string {
	return fmt.Sprintf("%s:%s:%d:%s", emailAddress, accountToUpdateKYC, requestID, n.approvedRequestConfig.Emails.RequestTokenSuffix)
}

func (n *ReviewedKYCRequestNotifier) buildRejectedKYCUniqueToken(emailAddress, accountToUpdateKYC string, kycSequence uint32, requestID uint64) string {
	return fmt.Sprintf("%s:%s:%d:%d:%s", emailAddress, accountToUpdateKYC, kycSequence, requestID, n.rejectedRequestConfig.Emails.RequestTokenSuffix)
}
