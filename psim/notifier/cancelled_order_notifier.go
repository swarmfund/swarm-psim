package notifier

import (
	"context"
	"fmt"
	"strconv"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/internal"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
)

type CancelledOrderNotifier struct {
	emailSender          EmailSender
	eventConfig          EventConfig
	saleConnector        SaleConnector
	transactionConnector TransactionConnector
	userConnector        UserConnector

	checkSaleStateResponses <-chan horizon.CheckSaleStateResponse
}

func (n *CancelledOrderNotifier) listenAndProcessCancelledOrders(ctx context.Context) error {
	for {
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
				continue
			}

			err = n.processCheckSaleStateOperation(ctx, *checkSaleStateOp)
			if err != nil {
				return errors.Wrap(err, "failed to process CheckSaleState operation")
			}
		}
	}
}

func (n *CancelledOrderNotifier) processCheckSaleStateOperation(ctx context.Context, checkSaleStateOperation horizon.CheckSaleState) error {
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
		return errors.New("transaction doesn't exist")
	}

	ledgerChanges := internal.LedgerChanges(tx)

	for _, change := range ledgerChanges {
		offer := n.getRemovedOffer(change)
		if offer == nil {
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

func (n *CancelledOrderNotifier) getRemovedOffer(change xdr.LedgerEntryChange) *xdr.LedgerKeyOffer {
	if change.Type != xdr.LedgerEntryChangeTypeRemoved {
		return nil
	}

	removedEntryKey := change.Removed

	if removedEntryKey.Type != xdr.LedgerEntryTypeOfferEntry {
		return nil
	}

	offer := removedEntryKey.MustOffer()

	return &offer
}

func (n *CancelledOrderNotifier) notifyAboutCancelledOrder(ctx context.Context, offer xdr.LedgerKeyOffer, saleID uint64) error {
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
	emailUniqueToken := n.buildOrderCancelledUniqueToken(emailAddress, uint64(offer.OfferId), saleID)

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

func (n *CancelledOrderNotifier) buildOrderCancelledUniqueToken(emailAddress string, offerID, saleID uint64) string {
	return fmt.Sprintf("%s:%d:%d:%s", emailAddress, offerID, saleID, n.eventConfig.Emails.RequestTokenSuffix)
}

func (n *CancelledOrderNotifier) getSale(saleID uint64) (*horizon.Sale, error) {
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
