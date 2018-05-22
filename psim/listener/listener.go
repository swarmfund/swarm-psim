package listener

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
)

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
	Account string
	Name    OutputEventName
}

func NewOutputEvent(Account string, Name OutputEventName) *OutputEvent {
	return &OutputEvent{Account, Name}
}

// AppendedBy returns array of the receiver and new event from arguments.
func (oe *OutputEvent) AppendedBy(Account string, Name OutputEventName) (outputEvents []OutputEvent) {
	outputEvents = append([]OutputEvent{*oe}, *NewOutputEvent(Account, Name))
	return
}

// Alone returns arrays of one element which is the receiver.
func (oe *OutputEvent) Alone() (outputEvents []OutputEvent) {
	outputEvents = []OutputEvent{*oe}
	return
}

type OutputEventName string

const (
	OutputEventNameKycCreated            OutputEventName = "kyc_created"
	OutputEventNameKycUpdated            OutputEventName = "kyc_updated"
	OutputEventNameKycRejected           OutputEventName = "kyc_rejected"
	OutputEventNameKycApproved           OutputEventName = "kyc_approved"
	OutputEventNameUserReferred          OutputEventName = "user_referred"
	OutputEventNameFundsWithdrawn        OutputEventName = "funds_withdrawn"
	OutputEventNamePaymentV2Received     OutputEventName = "payment_v2_received"
	OutputEventNamePaymentV2Sent         OutputEventName = "payment_v2_sent"
	OutputEventNamePaymentReceived       OutputEventName = "payment_received"
	OutputEventNamePaymentSent           OutputEventName = "payment_sent"
	OutputEventNameFundsDeposited        OutputEventName = "funds_deposited"
	OutputEventNameFundsInvested         OutputEventName = "funds_invested"
	OutputEventNameReferredUserPassedKyc OutputEventName = "referred_user_passed_kyc"
)

func NewListener(requestProvider RequestProvider, txPacketStream <-chan horizon.TXPacket, accountsProvider AccountProvider, logger *logan.Entry) *Listener {
	return &Listener{
		requestProvider: requestProvider,
		txPacketStream:  txPacketStream,
		accountProvider: accountsProvider,
		logger:          logger,
	}
}

// Listen takes TransactionEvents from TxStreamer and outputs them sequentially to outputEventsStream.
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
		outputEvents = NewOutputEvent(referrerAddress, OutputEventNameUserReferred).Alone()
	}
	return
}

func (l *Listener) handleWithdrawRequest(txSourceAccount xdr.AccountId) (outputEvents []OutputEvent) {
	outputEvents = NewOutputEvent(txSourceAccount.Address(), OutputEventNameFundsWithdrawn).Alone()
	return
}

func (l *Listener) handlePaymentV2(txSourceAccount xdr.AccountId, opResultBody *xdr.PaymentV2Result) (outputEvents []OutputEvent) {
	if opResultBody == nil {
		l.logger.Warn("received nil body for paymentV2 op")
		return
	}
	if opResultBody.PaymentV2Response == nil {
		l.logger.Warn("received nil paymentV2response for paymentV2 op")
		return
	}
	outputEvents = NewOutputEvent(txSourceAccount.Address(), OutputEventNamePaymentV2Sent).
		AppendedBy(opResultBody.PaymentV2Response.Destination.Address(), OutputEventNamePaymentV2Received)
	return
}

func (l *Listener) handlePayment(txSourceAccount xdr.AccountId, opResultBody *xdr.PaymentResult) (outputEvents []OutputEvent) {
	if opResultBody == nil {
		l.logger.Warn("received nil body for payment op")
		return
	}
	if opResultBody.PaymentResponse == nil {
		l.logger.Warn("received nil paymentResponse for payment op")
		return
	}
	outputEvents = NewOutputEvent(txSourceAccount.Address(), OutputEventNamePaymentSent).
		AppendedBy(opResultBody.PaymentResponse.Destination.Address(), OutputEventNamePaymentReceived)
	return
}

func (l *Listener) handleManageOfferOp(txSourceAccount xdr.AccountId, opBody *xdr.ManageOfferOp) (outputEvents []OutputEvent) {
	if opBody == nil {
		l.logger.Warn("receive nil body for manage offer op")
		return
	}
	if opBody.OrderBookId != 0 && opBody.Amount != 0 {
		outputEvents = NewOutputEvent(txSourceAccount.Address(), OutputEventNameFundsInvested).Alone()
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
		outputEvents = NewOutputEvent(opSuccess.Receiver.Address(), OutputEventNameFundsDeposited).Alone()
	}
	return
}

func (l *Listener) handleKYCCreateUpdateRequestOp(opBody *xdr.CreateUpdateKycRequestOp) (outputEvents []OutputEvent) {
	if opBody == nil {
		l.logger.Warn("receive nil KYC create update req op body")
		return
	}
	if opBody.RequestId == 0 {
		outputEvents = NewOutputEvent(opBody.UpdateKycRequestData.AccountToUpdateKyc.Address(), OutputEventNameKycCreated).Alone()
		return
	} // if op.RequestId != 0
	outputEvents = NewOutputEvent(opBody.UpdateKycRequestData.AccountToUpdateKyc.Address(), OutputEventNameKycUpdated).Alone()
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
	outputEvents = NewOutputEvent(sourceAccount.Address(), OutputEventNameFundsDeposited).Alone()
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
		return NewOutputEvent(kycRequestDetails.AccountToUpdateKYC, OutputEventNameKycRejected).Alone()
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
			outputEvents = NewOutputEvent(kycRequestDetails.AccountToUpdateKYC, OutputEventNameKycApproved).Alone()
		}
		account, err := l.accountProvider.ByAddress(kycRequestDetails.AccountToUpdateKYC)
		if err != nil {
			l.logger.WithError(err).Warn(err, "cannot get account by address")
		}
		if account.Referrer != "" {
			outputEvents = append(outputEvents, *NewOutputEvent(account.Referrer, OutputEventNameReferredUserPassedKyc))
		}
		return
	}
	return nil
}
