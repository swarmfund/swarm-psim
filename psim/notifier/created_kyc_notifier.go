package notifier

import (
	"gitlab.com/swarmfund/horizon-connector/v2"
	"context"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"fmt"
	"gitlab.com/swarmfund/go/xdr"
)

type CreatedKYCNotifier struct {
	emailSender          EmailSender
	emailsConfig         EmailsConfig
	transactionConnector TransactionConnector
	userConnector        UserConnector

	createKYCRequestOpResponses <-chan horizon.CreateKYCRequestOpResponse
}

func (n *CreatedKYCNotifier) listenAndProcessCreatedKYCRequests(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	case createKYCRequestOpResponse, ok := <-n.createKYCRequestOpResponses:
		if !ok {
			return nil
		}

		createKYCRequestOp, err := createKYCRequestOpResponse.Unwrap()
		if err != nil {
			return errors.Wrap(err, "CreateKYCRequestOpListener sent error")
		}

		err = n.processCreateKYCRequestOperation(ctx, *createKYCRequestOp)
		if err != nil {
			return errors.Wrap(err, "failed to process CreateKYCRequest operation")
		}

		return nil
	}
}

func (n *CreatedKYCNotifier) processCreateKYCRequestOperation(ctx context.Context, createKYCRequestOperation horizon.CreateKYCRequestOp) error {
	txID := createKYCRequestOperation.TransactionID
	tx, err := n.transactionConnector.TransactionByID(txID)
	if err != nil {
		return errors.Wrap(err, "failed to get transaction", logan.F{
			"transaction_id": txID,
		})
	}
	if tx == nil {
		// Transaction doesn't exist
		return nil
	}

	ledgerChanges := tx.LedgerChanges()

	for _, change := range ledgerChanges {
		if !n.isCreatedKYCRequest(change) {
			continue
		}

		err := n.notifyAboutCreatedKYCRequest(ctx, createKYCRequestOperation.AccountToUpdateKYC, createKYCRequestOperation.ID)
		if err != nil {
			return errors.Wrap(err, "failed to notify about created KYC request", logan.F{
				"account_to_update_kyc": createKYCRequestOperation.AccountToUpdateKYC,
				"transaction_id":        createKYCRequestOperation.TransactionID,
			})
		}
	}

	return nil
}

func (n *CreatedKYCNotifier) isCreatedKYCRequest(change xdr.LedgerEntryChange) bool {
	if change.Type != xdr.LedgerEntryChangeTypeCreated {
		return false
	}

	createdEntry := change.Created

	if createdEntry.Data.Type != xdr.LedgerEntryTypeReviewableRequest {
		return false
	}

	createdReviewableRequest := createdEntry.Data.ReviewableRequest

	if createdReviewableRequest.Body.Type != xdr.ReviewableRequestTypeUpdateKyc {
		return false
	}

	return true
}

func (n *CreatedKYCNotifier) notifyAboutCreatedKYCRequest(ctx context.Context, accountToUpdateKYC string, operationID string) error {
	user, err := n.userConnector.User(accountToUpdateKYC)
	if err != nil {
		return errors.Wrap(err, "failed to load user", logan.F{
			"account_id": accountToUpdateKYC,
		})
	}
	if user == nil {
		return nil
	}

	emailAddress := user.Attributes.Email
	emailUniqueToken := n.buildCreatedKYCUniqueToken(emailAddress, accountToUpdateKYC, operationID)

	data := struct {
		Link string
	}{
		Link: n.emailsConfig.TemplateLinkURL,
	}

	err = n.emailSender.SendEmail(ctx, emailAddress, emailUniqueToken, data)
	if err != nil {
		return errors.Wrap(err, "failed to send email")
	}

	return nil
}

func (n *CreatedKYCNotifier) buildCreatedKYCUniqueToken(emailAddress, accountToUpdateKYC, operationID string) string {
	return fmt.Sprintf("%s:%s:%s:%s", emailAddress, accountToUpdateKYC, operationID, n.emailsConfig.RequestTokenSuffix)
}
