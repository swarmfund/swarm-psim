package notifier

import (
	"context"
	"fmt"
	"strconv"

	"math/big"

	"github.com/ethereum/go-ethereum/log"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/regources"
)

type PaymentV2Notifier struct {
	log                  *logan.Entry
	emailSender          EmailSender
	eventConfig          EventConfig
	transactionConnector TransactionConnector
	userConnector        UserConnector

	paymentV2Responses <-chan horizon.PaymentV2OpResponse
}

type PaymentParticipantsFees struct {
	SenderTotalFee   regources.Amount
	ReceiverTotalFee regources.Amount
}

type PaymentNotificationData struct {
	Payer                 string
	Receiver              string
	Amount                regources.Amount
	PaymentAsset          string
	SenderFee             regources.Amount
	SenderFeeAsset        string
	ReceiverFee           regources.Amount
	ReceiverFeeAsset      string
	SenderPaysForReceiver bool
	Link                  string
}

func (n *PaymentV2Notifier) listenAndProcessPaymentV2(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Info("Context is canceled - stopping runner.")
			return nil
		case paymentV2Response, ok := <-n.paymentV2Responses:
			if !ok {
				log.Info("PaymentV2 responses channel closed - stopping runner.")
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

func (n *PaymentV2Notifier) processPaymentV2Operation(
	ctx context.Context,
	paymentV2Operation horizon.PaymentV2Op,
) error {
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

func (n *PaymentV2Notifier) notifyAboutPaymentV2(
	ctx context.Context,
	paymentV2Operation horizon.PaymentV2Op,
) error {
	// FIXME: if one of the participants doesn't have user - no one will receive payment V2 notification

	sender, err := n.userConnector.User(paymentV2Operation.From)
	if err != nil {
		return errors.Wrap(err, "failed to load sender", logan.F{
			"account_id": paymentV2Operation.From,
		})
	}
	if sender == nil {
		return nil
	}

	receiver, err := n.userConnector.User(paymentV2Operation.To)
	if err != nil {
		return errors.Wrap(err, "failed to load receiver", logan.F{
			"account_id": paymentV2Operation.To,
		})
	}
	if receiver == nil {
		return nil
	}

	senderEmailAddress := sender.Attributes.Email
	senderEmailUniqueToken := n.buildPaymentV2UniqueToken(senderEmailAddress, paymentV2Operation.PaymentID)

	receiverEmailAddress := receiver.Attributes.Email
	receiverEmailUniqueToken := n.buildPaymentV2UniqueToken(receiverEmailAddress, paymentV2Operation.PaymentID)

	paymentParticipantFees, err := getPaymentParticipantsFees(paymentV2Operation)
	if err != nil {
		return errors.Wrap(err, "failed to get payment participants fees")
	}

	data := PaymentNotificationData{
		Payer:                 senderEmailAddress,
		Receiver:              receiverEmailAddress,
		Amount:                paymentV2Operation.Amount,
		PaymentAsset:          paymentV2Operation.Asset,
		SenderFee:             paymentParticipantFees.SenderTotalFee,
		SenderFeeAsset:        paymentV2Operation.SourceFeeData.ActualPaymentFeeAssetCode,
		ReceiverFee:           paymentParticipantFees.ReceiverTotalFee,
		ReceiverFeeAsset:      paymentV2Operation.DestinationFeeData.ActualPaymentFeeAssetCode,
		SenderPaysForReceiver: paymentV2Operation.SourcePaysForDest,
		Link: n.eventConfig.Emails.TemplateLinkURL,
	}

	err = n.emailSender.SendEmail(ctx, senderEmailAddress, senderEmailUniqueToken, data)
	if err != nil {
		return errors.Wrap(err, "failed to send email to payment sender")
	}

	err = n.emailSender.SendEmail(ctx, receiverEmailAddress, receiverEmailUniqueToken, data)
	if err != nil {
		return errors.Wrap(err, "failed to send email to payment receiver")
	}

	return nil
}

func (n *PaymentV2Notifier) buildPaymentV2UniqueToken(emailAddress string, paymentID uint64) string {
	return fmt.Sprintf("%s:%d:%s", emailAddress, paymentID, n.eventConfig.Emails.RequestTokenSuffix)
}

func getPaymentParticipantsFees(paymentV2Operation horizon.PaymentV2Op) (*PaymentParticipantsFees, error) {
	senderFixedFee := big.NewInt(int64(paymentV2Operation.SourceFeeData.FixedFee))
	senderPaymentFee := big.NewInt(int64(paymentV2Operation.SourceFeeData.ActualPaymentFee))
	senderTotalFee := big.NewInt(0).Add(senderFixedFee, senderPaymentFee)
	if !senderTotalFee.IsInt64() {
		return nil, errors.From(errors.New("Sender total fee overflows int64"), logan.F{
			"payment_id": paymentV2Operation.PaymentID,
		})
	}

	receiverFixedFee := big.NewInt(int64(paymentV2Operation.DestinationFeeData.FixedFee))
	receiverPaymentFee := big.NewInt(int64(paymentV2Operation.DestinationFeeData.ActualPaymentFee))
	receiverTotalFee := big.NewInt(0).Add(receiverFixedFee, receiverPaymentFee)
	if !receiverTotalFee.IsInt64() {
		return nil, errors.From(errors.New("Receiver total fee overflows int64"), logan.F{
			"payment_id": paymentV2Operation.PaymentID,
		})
	}

	paymentParticipantsFees := PaymentParticipantsFees{
		SenderTotalFee:   regources.Amount(senderTotalFee.Int64()),
		ReceiverTotalFee: regources.Amount(receiverTotalFee.Int64()),
	}

	return &paymentParticipantsFees, nil
}
