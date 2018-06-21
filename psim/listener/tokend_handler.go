package listener

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/kyc"
	"gitlab.com/swarmfund/psim/psim/listener/internal"
	"gitlab.com/tokend/go/xdr"
)

// TokendHandler is a Handler implementation to be used with tokend stuff
type TokendHandler struct {
	processorByOpType map[xdr.OperationType]Processor
	HorizonConnector  *horizon.Connector
	logger            *logan.Entry
}

// NewTokendHandler constructs a TokendHandler without any binded processors
func NewTokendHandler(logger *logan.Entry, connector *horizon.Connector) *TokendHandler {
	return &TokendHandler{make(map[xdr.OperationType]Processor), connector, logger}
}

func (th *TokendHandler) withTokendProcessors() *TokendHandler {
	th.SetProcessor(xdr.OperationTypeCreateKycRequest, processKYCCreateUpdateRequestOp)
	th.SetProcessor(xdr.OperationTypeReviewRequest, processReviewRequestOp(th.HorizonConnector.Operations(), th.HorizonConnector.Accounts()))
	th.SetProcessor(xdr.OperationTypeCreateIssuanceRequest, processCreateIssuanceRequestOp)
	th.SetProcessor(xdr.OperationTypeManageOffer, processManageOfferOp)
	th.SetProcessor(xdr.OperationTypePayment, processPayment)
	th.SetProcessor(xdr.OperationTypePaymentV2, processPaymentV2)
	th.SetProcessor(xdr.OperationTypeCreateWithdrawalRequest, processWithdrawRequest)
	th.SetProcessor(xdr.OperationTypeCreateAccount, processCreateAccountOp)
	return th
}

// SetProcessor binds a processor to specified opType
func (th TokendHandler) SetProcessor(opType xdr.OperationType, processor Processor) {
	th.processorByOpType[opType] = processor
}

// UserData contains actor name and email to be sent to analytics
type UserData struct {
	Name    string
	Email   string
	Country string
}

func (th TokendHandler) lookupUserData(event *BroadcastedEvent) (*UserData, error) {
	user, err := th.HorizonConnector.Users().User(event.Account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup user by id")
	}
	if user == nil {
		return nil, errors.Wrap(err, "user not found")
	}
	account, err := th.HorizonConnector.Accounts().ByAddress(event.Account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup account by address")
	}
	if account == nil {
		return nil, errors.Wrap(err, "account not found")
	}

	blobs := th.HorizonConnector.Blobs()

	accountKycData := account.KYC.Data
	if accountKycData == nil {
		return nil, errors.New("nil account kyc data")
	}

	blob, err := blobs.Blob(accountKycData.BlobID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup blob by id")
	}
	if blob == nil {
		return nil, errors.Wrap(err, "blob not found")
	}

	kycData, err := kyc.ParseKYCData(blob.Attributes.Value)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse kyc data")
	}

	var name string
	if kycData != nil {
		name := kycData.FirstName + " " + kycData.LastName
	}

	return &UserData{name, user.Attributes.Email, kycData.Address.Country}, nil
}

// Process starts processing all op using processors in the map
func (th TokendHandler) Process(ctx context.Context, extractedItems <-chan ExtractedItem) <-chan ProcessedItem {
	broadcastedEvents := make(chan ProcessedItem)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				th.logger.WithRecover(r).Error("panic while processing event")
			}
			close(broadcastedEvents)
		}()

		for extractedItem := range extractedItems {
			if extractedItem.Error != nil {
				th.logger.WithError(extractedItem.Error).Warn("got invalid extracted item, skipping")
				continue
			}

			opData := extractedItem.ExtractedOpData
			opType := opData.Op.Body.Type

			process := th.processorByOpType[opType]
			if process == nil {
				continue
			}

			events := process(opData)

			for _, event := range events {
				select {
				case <-ctx.Done():
					return
				default:
				}

				if event.Error != nil {
					th.logger.WithError(event.Error).Warn("got invalid event, skipping")
					continue
				}
				if event.BroadcastedEvent == nil {
					th.logger.Info("got empty event")
					continue
				}

				userData, err := th.lookupUserData(event.BroadcastedEvent)
				if err != nil {
					th.logger.WithError(err).Warn("userdata lookup failed")
					continue
				}

				if userData == nil {
					th.logger.Info("no userdata, skipping")
					continue
				}

				broadcastedEvent := internal.ValidProcessedItem(event.BroadcastedEvent.WithActor(userData.Name, userData.Email))
				if broadcastedEvent.BroadcastedEvent.Name == BroadcastedEventNameFundsInvested {
					broadcastedEvent.BroadcastedEvent.InvestmentCountry = userData.Country
				}

				broadcastedEvents <- *broadcastedEvent
			}
		}
	}()

	return broadcastedEvents
}
