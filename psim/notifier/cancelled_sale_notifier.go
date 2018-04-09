package notifier

import (
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/go/xdr"
	"context"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"fmt"
	"strconv"
)

type CancelledSaleNotifier struct {
	emailSender          EmailSender
	eventConfig          EventConfig
	saleConnector        SaleConnector
	transactionConnector TransactionConnector
	userConnector        UserConnector

	checkSaleStateResponses <-chan horizon.CheckSaleStateResponse
}

func (n *CancelledSaleNotifier) listenAndProcessCancelledSales(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	case checkSaleStateResponse, ok := <-n.checkSaleStateResponses:
		if !ok {
			return nil
		}

		checkSaleStateOp, err := checkSaleStateResponse.Unwrap()
		if err != nil {
			return errors.Wrap(err, "CheckSaleStateListener sent error")
		}

		cursor, err := strconv.ParseUint(checkSaleStateOp.PT, 10, 64)
		if err != nil {
			return errors.Wrap(err, "failed to parse paging token", logan.F{
				"paging_token": checkSaleStateOp.PT,
			})
		}
		if cursor < n.eventConfig.Cursor {
			return nil
		}

		err = n.processCheckSaleStateOperation(ctx, *checkSaleStateOp)
		if err != nil {
			return errors.Wrap(err, "failed to process CheckSaleState operation")
		}

		return nil
	}
}

func (n *CancelledSaleNotifier) processCheckSaleStateOperation(ctx context.Context, checkSaleStateOperation horizon.CheckSaleState) error {
	if checkSaleStateOperation.Effect != xdr.CheckSaleStateEffectUpdated.String() {
		return nil
	}

	txID := checkSaleStateOperation.TransactionID
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
		offer, ok := n.getRemovedOffer(change)
		if !ok {
			continue
		}
		err = n.notifyAboutCancelledOrder(ctx, *offer, checkSaleStateOperation.SaleID)
		if err != nil {
			return errors.Wrap(err, "failed to notify about cancelled order", logan.F{
				"sale_id": checkSaleStateOperation.SaleID,
			})
		}
	}

	return nil
}

func (n *CancelledSaleNotifier) getRemovedOffer(change xdr.LedgerEntryChange) (offerPtr *xdr.LedgerKeyOffer, ok bool) {
	if change.Type != xdr.LedgerEntryChangeTypeRemoved {
		return nil, false
	}

	removedEntryKey := change.Removed

	if removedEntryKey.Type != xdr.LedgerEntryTypeOfferEntry {
		return nil, false
	}

	offer := removedEntryKey.MustOffer()

	return &offer, true
}

func (n *CancelledSaleNotifier) notifyAboutCancelledOrder(ctx context.Context, offer xdr.LedgerKeyOffer, saleID uint64) error {
	ownerID := offer.OwnerId.Address()

	user, err := n.userConnector.User(ownerID)
	if err != nil {
		return errors.Wrap(err, "failed to load user", logan.F{
			"account_id": ownerID,
		})
	}
	if user == nil {
		// User doesn't exist
		return nil
	}

	emailAddress := user.Attributes.Email
	emailUniqueToken := n.buildSaleCancelledUniqueToken(emailAddress, uint64(offer.OfferId), saleID)

	sale, err := n.getSale(saleID)
	if err != nil {
		return errors.Wrap(err, "failed to get sale", logan.F{
			"sale_id": saleID,
		})
	}
	if sale == nil {
		// sale doesn't exist
		return nil
	}

	data := struct {
		Fund string
		Link string
	}{
		Fund: sale.Name(),
		Link: n.eventConfig.Emails.TemplateLinkURL,
	}

	err = n.emailSender.SendEmail(ctx, emailAddress, emailUniqueToken, data)
	if err != nil {
		return errors.Wrap(err, "failed to send email")
	}

	return nil
}

func (n *CancelledSaleNotifier) buildSaleCancelledUniqueToken(emailAddress string, offerID, saleID uint64) string {
	return fmt.Sprintf("%s:%d:%d:%s", emailAddress, offerID, saleID, n.eventConfig.Emails.RequestTokenSuffix)
}

func (n *CancelledSaleNotifier) getSale(saleID uint64) (*horizon.Sale, error) {
	sale, err := n.saleConnector.SaleByID(saleID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load sale", logan.F{
			"sale_id": saleID,
		})
	}
	if sale == nil {
		// Sale doesn't exist
		return nil, nil
	}
	if sale.Name() == "" {
		return nil, errors.New(fmt.Sprintf("invalid sale name for id: %d", saleID))
	}

	return sale, nil
}
