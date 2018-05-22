package listener

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
)

// TODO edit Makefile before push
// TODO edit Gopkg.* files before push (custom horizon-connector branch is used)
// TODO cursor = now before push
// TODO edit psim config.yaml before push
// TODO edit README.md if targets or operations to listen are changed

type Listener struct {
	requestProvider RequestProvider
	txPacketStream  <-chan horizon.TXPacket
	accountProvider AccountProvider
	logger          *logan.Entry
}

type AccountProvider interface {
	ByAddress(string) (*horizon.Account, error)
}

type RequestProvider interface {
	GetRequestByID(requestID uint64) (*horizon.Request, error)
}

type OutputEvent struct {
	Account   string
	EventName string
}

func NewOutputEvent(Account string, EventName string) *OutputEvent {
	return &OutputEvent{Account, EventName}
}

// returns array from the receiver and new event from arguments
func (oe *OutputEvent) AppendedBy(Account string, EventName string) (outputEvents []OutputEvent) {
	outputEvents = append([]OutputEvent{*oe}, *NewOutputEvent(Account, EventName))
	return
}

type lol []OutputEvent

func (l lol) Add() lol {
	return append(l, *NewOutputEvent("", ""))
}

// returns array of one element which is the receiver
func (oe *OutputEvent) Alone() (outputEvents []OutputEvent) {
	outputEvents = []OutputEvent{*oe}
	return
}

// TODO rename and fix values
const (
	KycCreatedOutputEventName               = "kyc_created"
	KycUpdatedOutputEventName               = "kyc_updated"
	KycRejectedOutputEventName              = "kyc_rejected"
	KycApprovedOutputEventName              = "kyc_approved"
	IssuanceRequestFulfilledOutputEventName = "issuance_request_fulfilled"
	UserReferredOutputEventName             = "user_referred"
	WithdrawOutputEventName                 = "withdraw"
	PaymentV2ReceivedOutputEventName        = "paymentV2_received"
	PaymentV2SentOutputEventName            = "paymentV2_sent"
	PaymentReceivedOutputEventName          = "payment_received"
	PaymentSentOutputEventName              = "payment_sent"
	IssuanceRequestOutputEventName          = "issuance_request"
	ManageOfferOutputEventName              = "manage_offer_done"
	ReferralPassedKYCOutputEventName        = "referral_passed_kyc"
)

func NewListener(requestProvider RequestProvider, txPacketStream <-chan horizon.TXPacket, accountsProvider AccountProvider, logger *logan.Entry) *Listener {
	return &Listener{
		requestProvider: requestProvider,
		txPacketStream:  txPacketStream,
		accountProvider: accountsProvider,
		logger:          logger,
	}
}

// Main Listener function, takes TransactionEvents from TxStreamer
// and outputs them sequentially to an `outputEventsStream` channel
func (l *Listener) Listen(ctx context.Context) <-chan OutputEvent {
	outputEventsStream := make(chan OutputEvent)
	go func() {
		defer func() {
			close(outputEventsStream)
		}()
		for receivedTx := range l.txPacketStream {
			select {
			case <-ctx.Done():
				return
			default:
				break
			}
			running.RunSafely(ctx, "", func(ctx context.Context) error {
				receivedTxBody, err := receivedTx.Unwrap()
				if err != nil {
					l.logger.WithError(err).Error("Bad tx received")
				}
				outEvents := l.handleOps(ctx, receivedTxBody.Transaction)

				for _, event := range outEvents {
					outputEventsStream <- event
				}
				return nil
			})
		}
	}()
	return outputEventsStream
}

// handle-all-needed-operations func
func (l *Listener) handleOps(ctx context.Context, tx *horizon.Transaction) (outputEvents []OutputEvent) {
	if tx == nil {
		l.logger.Info("Received nil tx")
		return
	}

	txEnv := tx.Envelope().Tx
	txSourceAccount := txEnv.SourceAccount
	txLedgerChanges := tx.GroupedLedgerChanges()
	ops := txEnv.Operations
	opsResults := tx.Result().Result.MustResults()

	for currentOpIndex, currentOp := range ops {
		running.RunSafely(ctx, "", func(ctx context.Context) (err error) {
			opLedgerChanges := txLedgerChanges[currentOpIndex]
			sourceAccount := txSourceAccount

			opResultTr := opsResults[currentOpIndex].Tr
			if opResultTr == nil {
				l.logger.Warn("received nil op result body")
				return nil
			}
			if currentOp.SourceAccount != nil {
				sourceAccount = *currentOp.SourceAccount
			}
			outputEvents = l.handleOp(currentOp, sourceAccount, opLedgerChanges, *opResultTr)

			if outputEvents == nil {
				return nil
			}
			return nil
		})
	}
	return
}

func (l *Listener) handleOp(
	op xdr.Operation, sourceAccount xdr.AccountId, opLedgerChanges []xdr.LedgerEntryChange,
	opResult xdr.OperationResultTr,
) []OutputEvent {
	switch op.Body.Type {
	case xdr.OperationTypeCreateKycRequest:
		return l.handleKYCCreateUpdateRequestOp(op.Body.CreateUpdateKycRequestOp)
	case xdr.OperationTypeReviewRequest:
		return l.handleReviewRequestOp(sourceAccount, op.Body.ReviewRequestOp, opLedgerChanges)
	case xdr.OperationTypeCreateIssuanceRequest:
		return l.handleCreateIssuanceRequestOp(opResult)
	case xdr.OperationTypeManageOffer:
		return l.handleManageOfferOp(sourceAccount, op.Body.ManageOfferOp)
	case xdr.OperationTypePayment:
		return l.handlePayment(sourceAccount, opResult.PaymentResult)
	case xdr.OperationTypePaymentV2:
		return l.handlePaymentV2(sourceAccount, opResult.PaymentV2Result)
	case xdr.OperationTypeCreateWithdrawalRequest:
		return l.handleWithdrawRequest(sourceAccount)
	case xdr.OperationTypeCreateAccount:
		return l.handleCreateAccountOp(op.Body.CreateAccountOp)
	}
	return nil
}

func (l *Listener) handleCreateAccountOp(opBody *xdr.CreateAccountOp) (outputEvents []OutputEvent) {
	if opBody == nil {
		l.logger.Warn("received nil body for create account op")
		return
	}
	referrer := opBody.Referrer
	if referrer == nil {
		l.logger.Info("received nil referrer for create account op")
		return
	}
	referrerAddress := referrer.Address()
	if referrerAddress != "" {
		outputEvents = NewOutputEvent(referrerAddress, UserReferredOutputEventName).Alone()
	}
	return
}

func (l *Listener) handleWithdrawRequest(txSourceAccount xdr.AccountId) (outputEvents []OutputEvent) {
	outputEvents = NewOutputEvent(txSourceAccount.Address(), WithdrawOutputEventName).Alone()
	return
}

func (l *Listener) handlePaymentV2(txSourceAccount xdr.AccountId, opResultBody *xdr.PaymentV2Result) (outputEvents []OutputEvent) {
	// TODO what events to emit if body is empty?
	if opResultBody == nil {
		l.logger.Warn("received nil body for paymentV2 op")
		return
	}
	if opResultBody.PaymentV2Response == nil {

		l.logger.Warn("received nil paymentV2response for paymentV2 op")
		return
	}
	outputEvents = NewOutputEvent(txSourceAccount.Address(), PaymentV2SentOutputEventName).
		AppendedBy(opResultBody.PaymentV2Response.Destination.Address(), PaymentV2ReceivedOutputEventName)
	return
}

func (l *Listener) handlePayment(txSourceAccount xdr.AccountId, opResultBody *xdr.PaymentResult) (outputEvents []OutputEvent) {
	// TODO what events to emit if body is empty?
	if opResultBody == nil {
		l.logger.Warn("received nil body for payment op")
		return
	}
	if opResultBody.PaymentResponse == nil {
		l.logger.Warn("received nil paymentResponse for payment op")
		return
	}
	outputEvents = NewOutputEvent(txSourceAccount.Address(), PaymentSentOutputEventName).
		AppendedBy(opResultBody.PaymentResponse.Destination.Address(), PaymentReceivedOutputEventName)
	return
}

func (l *Listener) handleManageOfferOp(txSourceAccount xdr.AccountId, opBody *xdr.ManageOfferOp) (outputEvents []OutputEvent) {
	if opBody == nil {
		l.logger.Warn("receive nil body for manage offer op")
		return
	}
	if opBody.OrderBookId != 0 && opBody.Amount != 0 {
		outputEvents = NewOutputEvent(txSourceAccount.Address(), ManageOfferOutputEventName).Alone()
	}
	return
}

func (l *Listener) handleCreateIssuanceRequestOp(opResult xdr.OperationResultTr) (outputEvents []OutputEvent) {
	if opResult.CreateIssuanceRequestResult == nil {
		l.logger.Warn("received nil create issuance req result body (result tr)")
		return
	}
	opSuccess := opResult.CreateIssuanceRequestResult.Success
	if opSuccess == nil {
		l.logger.Warn("received nil create issuance req result success")
		return
	}
	if opSuccess.Fulfilled == true {
		outputEvents = NewOutputEvent(opSuccess.Receiver.Address(), IssuanceRequestFulfilledOutputEventName).Alone()
	}
	return
}

func (l *Listener) handleKYCCreateUpdateRequestOp(opBody *xdr.CreateUpdateKycRequestOp) (outputEvents []OutputEvent) {
	if opBody == nil {
		l.logger.Warn("receive nil KYC create update req op body")
		return
	}
	if opBody.RequestId == 0 {
		outputEvents = NewOutputEvent(opBody.UpdateKycRequestData.AccountToUpdateKyc.Address(), KycCreatedOutputEventName).Alone()
		return
	} // if op.RequestId != 0
	outputEvents = NewOutputEvent(opBody.UpdateKycRequestData.AccountToUpdateKyc.Address(), KycUpdatedOutputEventName).Alone()
	return
}

func (l *Listener) handleReviewRequestOp(sourceAccount xdr.AccountId, op *xdr.ReviewRequestOp, ledgerEntryChanges []xdr.LedgerEntryChange) []OutputEvent {
	switch op.RequestDetails.RequestType {
	case xdr.ReviewableRequestTypeUpdateKyc:
		return l.handleKYCReview(op, ledgerEntryChanges)
	case xdr.ReviewableRequestTypeIssuanceCreate:
		return l.handleIssuanceCreateReq(sourceAccount)
	}
	return nil
}

func (l *Listener) handleIssuanceCreateReq(sourceAccount xdr.AccountId) (outputEvents []OutputEvent) {
	outputEvents = NewOutputEvent(sourceAccount.Address(), IssuanceRequestOutputEventName).Alone()
	return
}

func (l *Listener) handleKYCReview(opBody *xdr.ReviewRequestOp, ledgerChanges []xdr.LedgerEntryChange) (outputEvents []OutputEvent) {
	if opBody == nil {
		l.logger.Warn("nil review request opBody body")
		return
	}
	request, err := l.requestProvider.GetRequestByID(uint64(opBody.RequestId))
	if err != nil {
		l.logger.Warn("failed to get request by id")
		return nil
	}
	if request == nil {
		return nil
	}

	kycRequestDetails := request.Details.KYC

	if opBody.Action == xdr.ReviewRequestOpActionReject || opBody.Action == xdr.ReviewRequestOpActionPermanentReject {
		return NewOutputEvent(kycRequestDetails.AccountToUpdateKYC, KycRejectedOutputEventName).Alone()
	}

	if opBody.Action != xdr.ReviewRequestOpActionApprove {
		return nil
	}

	for _, ledgerChange := range ledgerChanges {
		if ledgerChange.Removed == nil {
			return nil
		}
		if ledgerChange.Removed.ReviewableRequest == nil {
			return nil
		}
		reviewableRequest := ledgerChange.Removed.ReviewableRequest
		if opBody.RequestId == reviewableRequest.RequestId {
			outputEvents = NewOutputEvent(kycRequestDetails.AccountToUpdateKYC, KycApprovedOutputEventName).Alone()
		}
		account, err := l.accountProvider.ByAddress(kycRequestDetails.AccountToUpdateKYC)
		if err != nil {
			l.logger.WithError(err).Warn(err, "cannot get account by address")
		}
		if account.Referrer != "" {
			outputEvents = append(outputEvents, *NewOutputEvent(account.Referrer, ReferralPassedKYCOutputEventName))
		}
		return
	}
	return nil
}
