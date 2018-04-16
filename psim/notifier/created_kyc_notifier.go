package notifier

import (
	"context"
	"fmt"
	"strconv"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

type CreatedKYCNotifier struct {
	emailSender          EmailSender
	eventConfig          EventConfig
	transactionConnector TransactionConnector
	userConnector        UserConnector
	kycDataHelper        KYCDataHelper

	createKYCRequestOpResponses <-chan horizon.CreateKYCRequestOpResponse
}

func (n *CreatedKYCNotifier) listenAndProcessCreatedKYCRequests(ctx context.Context) error {
	for {
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
			fields := logan.F{
				"operation_id": createKYCRequestOp.ID,
				"request_id":   createKYCRequestOp.RequestID,
				"paging_token": createKYCRequestOp.PT,
				"tx_id":        createKYCRequestOp.TransactionID,
			}

			cursor, err := strconv.ParseUint(createKYCRequestOp.PT, 10, 64)
			if err != nil {
				return errors.Wrap(err, "failed to parse paging token", fields)
			}

			if !n.canNotifyAboutCreatedKYC(cursor) {
				continue
			}

			err = n.processCreateKYCRequestOperation(ctx, *createKYCRequestOp)
			if err != nil {
				return errors.Wrap(err, "Failed to process CreateKYCRequest operation", fields)
			}
		}
	}
}

func (n *CreatedKYCNotifier) canNotifyAboutCreatedKYC(cursor uint64) bool {
	return cursor >= n.eventConfig.Cursor
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
		return errors.New("transaction doesn't exist")
	}

	// we need ledger changes to ensure that KYCRequest was created but not updated through CreateKYCRequestOperation
	ledgerChanges := tx.LedgerChanges()

	for _, change := range ledgerChanges {
		if !n.isCreatedKYCRequest(change) {
			continue
		}

		err := n.notifyAboutCreatedKYCRequest(ctx, createKYCRequestOperation)
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

func (n *CreatedKYCNotifier) notifyAboutCreatedKYCRequest(ctx context.Context, createKYCRequestOperation horizon.CreateKYCRequestOp) error {
	if createKYCRequestOperation.AccountTypeToSet != int32(xdr.AccountTypeGeneral) {
		return nil
	}

	user, err := n.userConnector.User(createKYCRequestOperation.AccountToUpdateKYC)
	if err != nil {
		return errors.Wrap(err, "failed to load user", logan.F{
			"account_id": createKYCRequestOperation.AccountToUpdateKYC,
		})
	}
	if user == nil {
		return nil
	}

	blobKYCData, err := n.kycDataHelper.getBlobKYCData(createKYCRequestOperation.KYCData)
	if err != nil {
		return errors.Wrap(err, "failed to get blob KYC data")
	}

	emailAddress := user.Attributes.Email
	emailUniqueToken := n.buildCreatedKYCUniqueToken(emailAddress, createKYCRequestOperation.AccountToUpdateKYC, createKYCRequestOperation.RequestID)

	data := struct {
		Link      string
		FirstName string
	}{
		Link:      n.eventConfig.Emails.TemplateLinkURL,
		FirstName: blobKYCData.FirstName,
	}

	err = n.emailSender.SendEmail(ctx, emailAddress, emailUniqueToken, data)
	if err != nil {
		return errors.Wrap(err, "failed to send email")
	}

	return nil
}

func (n *CreatedKYCNotifier) buildCreatedKYCUniqueToken(emailAddress, accountToUpdateKYC string, requestID uint64) string {
	return fmt.Sprintf("%s:%s:%d:%s", emailAddress, accountToUpdateKYC, requestID, n.eventConfig.Emails.RequestTokenSuffix)
}
