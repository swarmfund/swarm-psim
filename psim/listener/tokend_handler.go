package listener

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/kyc"
	"gitlab.com/swarmfund/psim/psim/listener/internal"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
)

// TokendHandler is a Handler implementation to be used with tokend stuff
type TokendHandler struct {
	processorByOpType map[xdr.OperationType]Processor
	HorizonConnector  horizon.Connector
}

// NewTokendHandler constructs a TokendHandler without any binded processors
func NewTokendHandler() *TokendHandler {
	return &TokendHandler{make(map[xdr.OperationType]Processor), horizon.Connector{}}
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

func (th TokendHandler) lookupUserData(event MaybeBroadcastedEvent) (UserData, error) {
	user, err := th.HorizonConnector.Users().User(event.BroadcastedEvent.Account)
	if err != nil {
		return UserData{}, errors.Wrap(event.Error, "failed to lookup user by id")
	}
	if user != nil {
		return UserData{}, errors.Wrap(event.Error, "user not found")
	}

	account, err := th.HorizonConnector.Accounts().ByAddress(event.BroadcastedEvent.Account)
	if err != nil {
		return UserData{}, errors.Wrap(event.Error, "failed to lookup account by address")
	}
	if account != nil {
		return UserData{}, errors.Wrap(event.Error, "account not found")
	}

	blob, err := th.HorizonConnector.Blobs().Blob(account.KYC.Data.BlobID)
	if err != nil {
		return UserData{}, errors.Wrap(event.Error, "failed to lookup blob by id")
	}
	if blob != nil {
		return UserData{}, errors.Wrap(event.Error, "blob not found")
	}

	kycData, err := kyc.ParseKYCData(blob.Attributes.Value)
	if err != nil {
		return UserData{}, errors.Wrap(event.Error, "failed to parse kyc data")
	}

	name := kycData.FirstName + " " + kycData.LastName

	return UserData{name, user.Attributes.Email, kycData.Address.Country}, nil
}

// Process starts processing all op using processors in the map
func (th TokendHandler) Process(ctx context.Context, extractedItems <-chan ExtractedItem) <-chan ProcessedItem {
	broadcastedEvents := make(chan ProcessedItem)

	go func() {
		defer func() {
			close(broadcastedEvents)
		}()

		for extractedItem := range extractedItems {
			select {
			case <-ctx.Done():
				return
			default:
			}

			if extractedItem.Error != nil {
				broadcastedEvents <- *internal.InvalidProcessedItem(errors.Wrap(extractedItem.Error, "received invalid extracted item"))
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
				if event.Error != nil {
					broadcastedEvents <- *internal.InvalidProcessedItem(errors.Wrap(event.Error, "failed to process op data from extracted item"))
					continue
				}

				userData, err := th.lookupUserData(event)
				if err != nil {
					broadcastedEvents <- *internal.InvalidProcessedItem(errors.Wrap(err, "failed to lookup user data"))
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
