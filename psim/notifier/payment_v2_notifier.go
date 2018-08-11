package notifier

import (
	"gitlab.com/tokend/horizon-connector"
	"context"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"strconv"
	"fmt"
	"gitlab.com/tokend/regources"
)

type PaymentV2Notifier struct {
	emailSender          EmailSender
	eventConfig          EventConfig
	transactionConnector TransactionConnector
	userConnector        UserConnector

	paymentV2Responses <-chan horizon.PaymentV2OpResponse
}

func (n *PaymentV2Notifier) listenAndProcessPaymentV2(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case paymentV2Response, ok := <-n.paymentV2Responses:
			if !ok {
				return nil
			}

			paymentV2Op, err := paymentV2Response.Unwrap()
			if err != nil {
				return errors.Wrap(err, "payment op v2 listener sent error")
			}
			fields := logan.F{
				"tx_id":        paymentV2Op.TransactionID,
				"paging_token": paymentV2Op.PT,
			}

			cursor, err := strconv.ParseUint(paymentV2Op.PT, 10, 64)
			if err != nil {
				return errors.Wrap(err, "failed to parse paging token", fields)
			}

			if !n.canNotifyAboutPaymentV2(cursor) {
				continue
			}

			err = n.processPaymentV2Operation(ctx, *paymentV2Op)
			if err != nil {
				return errors.Wrap(err, "failed to process PaymentV2 operation", fields)
			}
		}
	}
}

func (n *PaymentV2Notifier) canNotifyAboutPaymentV2(cursor uint64) bool {
	return cursor >= n.eventConfig.Cursor
}

func (n *PaymentV2Notifier) processPaymentV2Operation(ctx context.Context,
	paymentV2Operation horizon.PaymentV2Op) error {
	err := n.notifyAboutPaymentV2(ctx, paymentV2Operation)
	if err != nil {
		return errors.Wrap(err, "failed to notify about payment V2",
			logan.F{
				"payment_id":     paymentV2Operation.PaymentID,
				"destination":    paymentV2Operation.To,
				"amount":         paymentV2Operation.Amount,
				"transaction_id": paymentV2Operation.TransactionID,
			})
	}

	return nil
}

func (n *PaymentV2Notifier) notifyAboutPaymentV2(ctx context.Context, paymentV2Operation horizon.PaymentV2Op) error {
	user, err := n.userConnector.User(paymentV2Operation.To)
	if err != nil {
		return errors.Wrap(err, "failed to load user", logan.F{
			"account_id": paymentV2Operation.To,
		})
	}
	if user == nil {
		return nil
	}

	emailAddress := user.Attributes.Email
	emailUniqueToken := n.buildPaymentV2UniqueToken(emailAddress, paymentV2Operation.PaymentID)

	data := struct {
		Payer                 string
		Receiver              string
		Amount                regources.Amount
		PaymentAsset          string
		SenderFee             regources.Amount
		SenderFeeAsset        string
		ReceiverFee           regources.Amount
		ReceiverFeeAsset      string
		SenderPaysForReceiver bool
	}{
		Payer:                 paymentV2Operation.From,
		Receiver:              paymentV2Operation.To,
		Amount:                paymentV2Operation.Amount,
		PaymentAsset:          paymentV2Operation.Asset,
		SenderFee:             paymentV2Operation.SourceFeeData.FixedFee + paymentV2Operation.SourceFeeData.ActualPaymentFee,
		SenderFeeAsset:        paymentV2Operation.SourceFeeData.ActualPaymentFeeAssetCode,
		ReceiverFee:           paymentV2Operation.DestinationFeeData.FixedFee + paymentV2Operation.DestinationFeeData.ActualPaymentFee,
		ReceiverFeeAsset:      paymentV2Operation.DestinationFeeData.ActualPaymentFeeAssetCode,
		SenderPaysForReceiver: paymentV2Operation.SourcePaysForDest,
	}

	err = n.emailSender.SendEmail(ctx, emailAddress, emailUniqueToken, data)
	if err != nil {
		return errors.Wrap(err, "failed to send email")
	}

	return nil
}

func (n *PaymentV2Notifier) buildPaymentV2UniqueToken(emailAddress string, paymentID uint64) string {
	return fmt.Sprintf("%s:%d:%s", emailAddress, paymentID, n.eventConfig.Emails.RequestTokenSuffix)
}
