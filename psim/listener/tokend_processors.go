package listener

import (
	"gitlab.com/swarmfund/psim/psim/listener/internal"

	"time"

	"gitlab.com/tokend/go/xdr"
)

func (th *TokendHandler) processCreateAccountOp(txData TxData) (outputEvents []BroadcastedEvent) {
	opBody := txData.Op.Body.CreateAccountOp
	if opBody == nil {
		return
	}
	referrer := opBody.Referrer
	if referrer == nil {
		return
	}
	referrerAddress := referrer.Address()
	if referrerAddress != "" {
		outputEvents = internal.NewBroadcastedEvent(referrerAddress, OutputEventNameUserReferred, txData.CreatedAt).Alone()
	}
	return outputEvents
}

func (th *TokendHandler) processWithdrawRequest(txData TxData) (outputEvents []BroadcastedEvent) {
	txSourceAccount := txData.SourceAccount
	outputEvents = internal.NewBroadcastedEvent(txSourceAccount.Address(), OutputEventNameFundsWithdrawn, txData.CreatedAt).Alone()
	return
}

func (th *TokendHandler) processPaymentV2(txData TxData) (outputEvents []BroadcastedEvent) {
	txSourceAccount := txData.SourceAccount
	opResultBody := txData.OpResult.PaymentV2Result
	if opResultBody == nil {
		return
	}
	if opResultBody.PaymentV2Response == nil {
		return
	}
	outputEvents = internal.NewBroadcastedEvent(txSourceAccount.Address(), OutputEventNamePaymentV2Sent, txData.CreatedAt).
		AppendedBy(opResultBody.PaymentV2Response.Destination.Address(), OutputEventNamePaymentV2Received, txData.CreatedAt)
	return
}

func (th *TokendHandler) processPayment(txData TxData) (outputEvents []BroadcastedEvent) {
	txSourceAccount := txData.SourceAccount
	opResultBody := txData.OpResult.PaymentResult
	if opResultBody == nil {
		return
	}
	if opResultBody.PaymentResponse == nil {
		return
	}
	outputEvents = internal.NewBroadcastedEvent(txSourceAccount.Address(), OutputEventNamePaymentSent, txData.CreatedAt).
		AppendedBy(opResultBody.PaymentResponse.Destination.Address(), OutputEventNamePaymentReceived, txData.CreatedAt)
	return
}

func (th *TokendHandler) processManageOfferOp(txData TxData) (outputEvents []BroadcastedEvent) {
	txSourceAccount := txData.SourceAccount
	opBody := txData.Op.Body.ManageOfferOp
	if opBody == nil {
		return
	}
	if opBody.OrderBookId != 0 && opBody.Amount != 0 {
		outputEvents = internal.NewBroadcastedEvent(txSourceAccount.Address(), OutputEventNameFundsInvested, txData.CreatedAt).Alone()
	}
	return
}

func (th *TokendHandler) processCreateIssuanceRequestOp(txData TxData) (outputEvents []BroadcastedEvent) {
	opResult := txData.OpResult
	if opResult.CreateIssuanceRequestResult == nil {
		return
	}
	opSuccess := opResult.CreateIssuanceRequestResult.Success
	if opSuccess == nil {
		return
	}
	if opSuccess.Fulfilled == true {
		outputEvents = internal.NewBroadcastedEvent(opSuccess.Receiver.Address(), OutputEventNameFundsDeposited, txData.CreatedAt).Alone()
	}
	return
}

func (th *TokendHandler) processKYCCreateUpdateRequestOp(txData TxData) (outputEvents []BroadcastedEvent) {
	opBody := txData.Op.Body.CreateUpdateKycRequestOp
	if opBody == nil {
		return
	}
	if opBody.RequestId == 0 {
		outputEvents = internal.NewBroadcastedEvent(opBody.UpdateKycRequestData.AccountToUpdateKyc.Address(), OutputEventNameKycCreated, txData.CreatedAt).Alone()
		return
	} // if op.RequestId != 0
	outputEvents = internal.NewBroadcastedEvent(opBody.UpdateKycRequestData.AccountToUpdateKyc.Address(), OutputEventNameKycUpdated, txData.CreatedAt).Alone()
	return
}

// TODO generalize this method, too
func (th *TokendHandler) processReviewRequestOp(txData TxData) (outputEvents []BroadcastedEvent) {

	sourceAccount := txData.SourceAccount
	op := txData.Op.Body.ReviewRequestOp
	ledgerEntryChanges := txData.OpLedgerChanges
	time := txData.CreatedAt

	var err error
	switch op.RequestDetails.RequestType {
	case xdr.ReviewableRequestTypeUpdateKyc:
		outputEvents, err = th.handleKYCReview(op, ledgerEntryChanges, time)
	case xdr.ReviewableRequestTypeIssuanceCreate:
		outputEvents = th.handleIssuanceCreateReq(sourceAccount, time)
	}
	if err != nil {
	}
	return outputEvents
}

func (th *TokendHandler) handleIssuanceCreateReq(sourceAccount xdr.AccountId, time *time.Time) []BroadcastedEvent {
	return internal.NewBroadcastedEvent(sourceAccount.Address(), OutputEventNameFundsDeposited, time).Alone()
}

func (th *TokendHandler) handleKYCReview(opBody *xdr.ReviewRequestOp, ledgerChanges []xdr.LedgerEntryChange, time *time.Time) (outputEvents []BroadcastedEvent, errx error) {
	if opBody == nil {
		return nil, nil
	}
	request, err := th.requestsProvider.GetRequestByID(uint64(opBody.RequestId))
	if err != nil {
		return nil, nil
	}

	if request == nil {
		return nil, nil
	}

	kycRequestDetails := request.Details.KYC

	if opBody.Action == xdr.ReviewRequestOpActionReject || opBody.Action == xdr.ReviewRequestOpActionPermanentReject {
		return internal.NewBroadcastedEvent(kycRequestDetails.AccountToUpdateKYC, OutputEventNameKycRejected, time).Alone(), nil
	}

	if opBody.Action != xdr.ReviewRequestOpActionApprove {
		return nil, nil
	}

	for _, ledgerChange := range ledgerChanges {
		if ledgerChange.Removed == nil {
			continue
		}
		if ledgerChange.Removed.ReviewableRequest == nil {
			continue
		}
		reviewableRequest := ledgerChange.Removed.ReviewableRequest
		if opBody.RequestId == reviewableRequest.RequestId {
			outputEvents = internal.NewBroadcastedEvent(kycRequestDetails.AccountToUpdateKYC, OutputEventNameKycApproved, time).Alone()
		}
		account, err := th.accountsProvider.ByAddress(kycRequestDetails.AccountToUpdateKYC)
		if err != nil {
		}
		if account.Referrer != "" {
			outputEvents = append(outputEvents, *internal.NewBroadcastedEvent(account.Referrer, OutputEventNameReferredUserPassedKyc, time))
		}
	}
	return
}
