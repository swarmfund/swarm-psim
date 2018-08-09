package eventsubmitter

import (
	"strings"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/swarmfund/psim/psim/eventsubmitter/internal"

	horizon "gitlab.com/tokend/horizon-connector"

	"gitlab.com/tokend/go/xdr"
)

// AccountProvider is responsible for account lookup by address
type AccountProvider interface {
	ByAddress(string) (*horizon.Account, error)
}

// RequestProvider is responsible for request lookup by address
type RequestProvider interface {
	GetRequestByID(requestID uint64) (*horizon.Request, error)
}

// Event names sent to analytics services
const (
	BroadcastedEventNameKycCreated            BroadcastedEventName = "kyc_created"
	BroadcastedEventNameKycUpdated            BroadcastedEventName = "kyc_updated"
	BroadcastedEventNameKycRejected           BroadcastedEventName = "kyc_rejected"
	BroadcastedEventNameKycApproved           BroadcastedEventName = "kyc_approved"
	BroadcastedEventNameUserReferred          BroadcastedEventName = "user_referred"
	BroadcastedEventNameFundsWithdrawn        BroadcastedEventName = "funds_withdrawn"
	BroadcastedEventNamePaymentV2Received     BroadcastedEventName = "payment_v2_received"
	BroadcastedEventNamePaymentV2Sent         BroadcastedEventName = "payment_v2_sent"
	BroadcastedEventNamePaymentReceived       BroadcastedEventName = "payment_received"
	BroadcastedEventNamePaymentSent           BroadcastedEventName = "payment_sent"
	BroadcastedEventNameFundsDeposited        BroadcastedEventName = "funds_deposited"
	BroadcastedEventNameFundsInvested         BroadcastedEventName = "funds_invested"
	BroadcastedEventNameReferredUserPassedKyc BroadcastedEventName = "referred_user_passed_kyc"
	BroadcastedEventNameReceivedAirdrop       BroadcastedEventName = "received_airdrop"
)

func processCreateAccountOp(opData OpData) []MaybeBroadcastedEvent {
	opBody := opData.Op.Body.CreateAccountOp
	if opBody == nil {
		return internal.InvalidBroadcastedEvent(errors.New("received nil create account op body")).Alone()
	}

	referrer := opBody.Referrer

	if referrer == nil {
		return nil
	}

	referrerAddress := referrer.Address()

	if referrerAddress == "" {
		return nil
	}

	return internal.ValidBroadcastedEvent(referrerAddress, BroadcastedEventNameUserReferred, opData.CreatedAt).Alone()
}

func processWithdrawRequest(opData OpData) []MaybeBroadcastedEvent {
	txSourceAccountAddress := opData.SourceAccount.Address()

	return internal.ValidBroadcastedEvent(txSourceAccountAddress, BroadcastedEventNameFundsWithdrawn, opData.CreatedAt).Alone()
}

func processPaymentV2(opData OpData) []MaybeBroadcastedEvent {
	txSourceAccount := opData.SourceAccount
	opResultBody := opData.OpResult.PaymentV2Result

	if opResultBody == nil {
		return internal.InvalidBroadcastedEvent(errors.New("received nil payment v2 op result body")).Alone()
	}

	if opResultBody.PaymentV2Response == nil {
		return internal.InvalidBroadcastedEvent(errors.New("received nil payment v2 response")).Alone()
	}

	outputEvents := internal.ValidBroadcastedEvent(txSourceAccount.Address(), BroadcastedEventNamePaymentV2Sent, opData.CreatedAt).
		AppendedBy(opResultBody.PaymentV2Response.Destination.Address(), BroadcastedEventNamePaymentV2Received, opData.CreatedAt)

	return outputEvents
}

func processPayment(opData OpData) []MaybeBroadcastedEvent {
	txSourceAccount := opData.SourceAccount
	opResultBody := opData.OpResult.PaymentResult

	if opResultBody == nil {
		return internal.InvalidBroadcastedEvent(errors.New("received nil payment op result body")).Alone()
	}

	if opResultBody.PaymentResponse == nil {
		return internal.InvalidBroadcastedEvent(errors.New("received nil payment response")).Alone()
	}

	outputEvents := internal.ValidBroadcastedEvent(txSourceAccount.Address(), BroadcastedEventNamePaymentSent, opData.CreatedAt).
		AppendedBy(opResultBody.PaymentResponse.Destination.Address(), BroadcastedEventNamePaymentReceived, opData.CreatedAt)

	return outputEvents
}

func processManageOfferOp(opData OpData) []MaybeBroadcastedEvent {
	txSourceAccount := opData.SourceAccount
	opBody := opData.Op.Body.ManageOfferOp

	if opBody == nil {
		return internal.InvalidBroadcastedEvent(errors.New("received nil manage offer op body")).Alone()
	}

	if opBody.OrderBookId == 0 {
		return nil
	}

	if opBody.Amount == 0 {
		return nil
	}

	validEvent := internal.ValidBroadcastedEvent(txSourceAccount.Address(), BroadcastedEventNameFundsInvested, opData.CreatedAt)
	validEvent.BroadcastedEvent.InvestmentAmount = int64(opBody.Amount)

	return validEvent.Alone()
}

func processCreateIssuanceRequestOp(opData OpData) []MaybeBroadcastedEvent {
	opBody := opData.Op.Body.CreateIssuanceRequestOp
	if opBody != nil {
		reference := opBody.Reference
		for _, airdropSuffix := range airdrop.AllAirdropSuffixes {
			if strings.HasSuffix(string(reference), airdropSuffix) {
				return internal.ValidBroadcastedEvent(opData.SourceAccount.Address(), BroadcastedEventNameReceivedAirdrop, opData.CreatedAt).Alone()
			}
		}
	}

	opResult := opData.OpResult

	if opResult.CreateIssuanceRequestResult == nil {
		return internal.InvalidBroadcastedEvent(errors.New("received nil create issuance req op body")).Alone()
	}

	opSuccess := opResult.CreateIssuanceRequestResult.Success

	if opSuccess == nil {
		return nil
	}

	if !opSuccess.Fulfilled {
		return nil
	}

	return internal.ValidBroadcastedEvent(opSuccess.Receiver.Address(), BroadcastedEventNameFundsDeposited, opData.CreatedAt).Alone()
}

func processKYCCreateUpdateRequestOp(opData OpData) []MaybeBroadcastedEvent {
	opBody := opData.Op.Body.CreateUpdateKycRequestOp

	if opBody == nil {
		return internal.InvalidBroadcastedEvent(errors.New("received nil op body")).Alone()
	}

	if opBody.RequestId == 0 {
		// kyc_created
		return internal.ValidBroadcastedEvent(opBody.UpdateKycRequestData.AccountToUpdateKyc.Address(), BroadcastedEventNameKycCreated, opData.CreatedAt).Alone()
	}

	// kyc_updated
	return internal.ValidBroadcastedEvent(opBody.UpdateKycRequestData.AccountToUpdateKyc.Address(), BroadcastedEventNameKycUpdated, opData.CreatedAt).Alone()
}

func processReviewRequestOp(requestsProvider RequestProvider, accountsProvider AccountProvider) Processor {
	return func(opData OpData) []MaybeBroadcastedEvent {
		requestType := opData.Op.Body.ReviewRequestOp.RequestDetails.RequestType

		switch requestType {
		case xdr.ReviewableRequestTypeUpdateKyc:
			return handleKYCReview(opData, requestsProvider, accountsProvider)
		case xdr.ReviewableRequestTypeIssuanceCreate:
			return handleIssuanceCreateReq(opData)
		}

		return nil
	}
}

func handleIssuanceCreateReq(opData OpData) []MaybeBroadcastedEvent {
	sourceAccountAddress := opData.SourceAccount.Address()
	time := opData.CreatedAt
	return internal.ValidBroadcastedEvent(sourceAccountAddress, BroadcastedEventNameFundsDeposited, time).Alone()
}

func findRemoval(ledgerChanges []xdr.LedgerEntryChange) *xdr.LedgerKeyReviewableRequest {
	for _, ledgerChange := range ledgerChanges {
		if ledgerChange.Removed == nil {
			continue
		}

		if ledgerChange.Removed.ReviewableRequest == nil {
			continue
		}

		return ledgerChange.Removed.ReviewableRequest
	}
	return nil
}

func getRequestKYCDetailsByID(id xdr.Uint64, requestsProvider RequestProvider) (*horizon.RequestKYCDetails, error) {
	request, err := requestsProvider.GetRequestByID(uint64(id))

	if err != nil {
		return nil, errors.Wrap(err, "failed to get request by id")
	}

	if request == nil {
		return nil, nil
	}

	return request.Details.KYC, nil
}

func findReferrer(accountsProvider AccountProvider, kycRequestDetails *horizon.RequestKYCDetails) (string, error) {
	account, err := accountsProvider.ByAddress(kycRequestDetails.AccountToUpdateKYC)

	if err != nil {
		return "", errors.New("failed to find account by address")
	}

	if account == nil {
		return "", errors.New("account by address doesn't exist")
	}

	if account.Referrer == "" {
		return "", nil
	}

	return account.Referrer, nil
}

func handleKYCReview(opData OpData, requestsProvider RequestProvider, accountsProvider AccountProvider) []MaybeBroadcastedEvent {
	ledgerChanges := opData.OpLedgerChanges
	time := opData.CreatedAt
	opBody := opData.Op.Body.ReviewRequestOp
	if opBody == nil {
		return internal.InvalidBroadcastedEvent(errors.New("received nil kyc review request body")).Alone()
	}

	kycRequestDetails, err := getRequestKYCDetailsByID(opBody.RequestId, requestsProvider)
	if err != nil {
		return internal.InvalidBroadcastedEvent(errors.Wrap(err, "failed to get kyc details")).Alone()
	}
	if kycRequestDetails == nil {
		return internal.InvalidBroadcastedEvent(errors.Wrap(err, "kyc request details not found")).Alone()
	}

	if kycRequestDetails.AccountToUpdateKYC == "" {
		return nil
	}

	if opBody.Action == xdr.ReviewRequestOpActionReject || opBody.Action == xdr.ReviewRequestOpActionPermanentReject {
		return internal.ValidBroadcastedEvent(kycRequestDetails.AccountToUpdateKYC, BroadcastedEventNameKycRejected, time).Alone()
	}

	if opBody.Action != xdr.ReviewRequestOpActionApprove {
		return internal.InvalidBroadcastedEvent(errors.New("kyc review request is neither approved nor rejected")).Alone()
	}

	var outputEvents []MaybeBroadcastedEvent

	reviewableRequest := findRemoval(ledgerChanges)
	if reviewableRequest == nil {
		return nil
	}
	if opBody.RequestId == reviewableRequest.RequestId {
		outputEvents = append(outputEvents, *internal.ValidBroadcastedEvent(kycRequestDetails.AccountToUpdateKYC, BroadcastedEventNameKycApproved, time))
	}

	referrer, err := findReferrer(accountsProvider, kycRequestDetails)
	if err != nil {
		outputEvents = append(outputEvents, *internal.InvalidBroadcastedEvent(errors.Wrap(err, "failed to find referrer")))
	}

	if referrer != "" {
		outputEvents = append(outputEvents, *internal.ValidBroadcastedEvent(referrer, BroadcastedEventNameReferredUserPassedKyc, time))
	}

	return outputEvents
}
