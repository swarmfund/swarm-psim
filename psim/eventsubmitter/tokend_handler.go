package eventsubmitter

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/eventsubmitter/internal"
	"gitlab.com/swarmfund/psim/psim/kyc"
	"gitlab.com/tokend/go/xdr"
	horizon "gitlab.com/tokend/horizon-connector"
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
	Name     string
	Email    string
	Country  string
	Referrer string
}

// TODO: consider use kycapi
func (th TokendHandler) lookupUserData(event *BroadcastedEvent) (*UserData, error) {
	user, err := th.HorizonConnector.Users().User(event.Account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup user by id", logan.F{
			"event_account": event.Account,
		})
	}
	if user == nil {
		return nil, nil
	}
	account, err := th.HorizonConnector.Accounts().ByAddress(event.Account)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup account by address", logan.F{
			"event_account": event.Account,
		})
	}
	if account == nil {
		return nil, nil
	}

	blobs := th.HorizonConnector.Blobs()

	accountKycData := account.KYC.Data
	var kycData *kyc.Data
	if accountKycData != nil {
		blob, err := blobs.Blob(accountKycData.BlobID)
		// returns are omitted intentionally to make events event without kyc-data
		if err != nil {
			th.logger.WithError(errors.Wrap(err, "failed to lookup blob by id", logan.F{
				"account_kyc_blob_id": accountKycData.BlobID,
			})).Error("failed to get blob from horizon")
		}
		if blob != nil {
			kycData, err = kyc.ParseKYCData(blob.Attributes.Value)
			if err != nil {
				th.logger.WithError(errors.Wrap(err, "failed to parse kyc data", logan.F{
					"kyc_attributes": blob.Attributes.Value,
				})).Error("got event with old kyc")
			}
		}
	}

	var name string
	var country string
	if kycData != nil {
		name = kycData.FirstName + " " + kycData.LastName
		country = kycData.Address.Country
	}

	return &UserData{
		Name:     name,
		Email:    user.Attributes.Email,
		Country:  country,
		Referrer: account.Referrer,
	}, nil
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
				th.logger.Debug("no suitable event processor")
				continue
			}

			events := process(opData)

			if events == nil {
				th.logger.Debug("no events produced after processing")
				continue
			}

			for _, event := range events {
				select {
				case <-ctx.Done():
					return
				default:
				}

				if event.Error != nil {
					th.logger.WithError(event.Error).WithFields(logan.F{
						"tx_paging_token": extractedItem.ExtractedOpData.PagingToken,
						"op_type":         opType,
					}).Warn("got invalid event, skipping")
					continue
				}
				if event.BroadcastedEvent == nil {
					th.logger.Debug("got empty event")
					continue
				}

				userData, err := th.lookupUserData(event.BroadcastedEvent)
				if err != nil {
					th.logger.WithError(err).WithFields(logan.F{
						"tx_paging_token": extractedItem.ExtractedOpData.PagingToken,
						"op_type":         opType,
					}).Error("userdata lookup failed")
					continue
				}

				var broadcastedEvent *ProcessedItem
				if userData == nil {
					th.logger.Debug("no user data found")
					broadcastedEvent = internal.ValidProcessedItem(event.BroadcastedEvent)
				} else {
					broadcastedEvent = internal.ValidProcessedItem(event.BroadcastedEvent.WithActor(userData.Name, userData.Email))
					broadcastedEvent.BroadcastedEvent.Country = userData.Country
				}
				broadcastedEvents <- *broadcastedEvent
			}
		}
	}()

	return broadcastedEvents
}
