package ordernotifier

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/notificator-server/client"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/app"
	"golang.org/x/net/context"
	"time"
	"fmt"
)

// UserConnector is an interface for retrieving specific user
type UserConnector interface {
	// User retrieves a single User by AccountID.
	// If User doesn't exist - nil,nil is returned.
	User(accountID string) (*horizon.User, error)
}

// TransactionConnector is an interface for retrieving transaction
// specified by provided transaction ID
type TransactionConnector interface {
	// TransactionByID retrieves Transaction with given transaction ID
	// If Transaction doesn't exist - nil,nil is returned.
	TransactionByID(txID string) (*horizon.Transaction, error)
}

// SaleConnector is an interface for retrieving sale
// specified by provided sale ID
type SaleConnector interface {
	// SaleByID retrieves Sale with given sale ID
	// If Sale doesn't exist - nil,nil is returned.
	SaleByID(saleID uint64) (*horizon.Sale, error)
}

type NotificatorConnector interface {
	Send(requestType int, token string, payload notificator.Payload) (*notificator.Response, error)
}

type Service struct {
	config               Config
	transactionConnector TransactionConnector
	emailSender          NotificatorConnector
	logger               *logan.Entry
	userConnector        UserConnector
	saleConnector        SaleConnector

	checkSaleStateResponses <-chan horizon.CheckSaleStateResponse
}

// New is a constructor of a service
func New(
	config Config,
	transactionConnector TransactionConnector,
	emailSender NotificatorConnector,
	logger *logan.Entry,
	userConnector UserConnector,
	saleConnector SaleConnector,
	checkSaleStateResponses <-chan horizon.CheckSaleStateResponse,
) *Service {
	return &Service{
		config:               config,
		transactionConnector: transactionConnector,
		emailSender:          emailSender,
		logger:               logger,
		userConnector:        userConnector,
		saleConnector:        saleConnector,

		checkSaleStateResponses: checkSaleStateResponses,
	}
}

func (s *Service) Run(ctx context.Context) {
	s.logger.Info("Starting...")
	app.RunOverIncrementalTimer(ctx, s.logger, "check_sale_state_operations_processor", s.listenAndProcessCheckSaleStateOperations, 0, 5*time.Second)
}

func (s *Service) listenAndProcessCheckSaleStateOperations(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	case checkSaleStateResponse, ok := <-s.checkSaleStateResponses:
		if !ok {
			return nil
		}

		checkSaleStateOp, err := checkSaleStateResponse.Unwrap()
		if err != nil {
			return errors.Wrap(err, "CheckSaleStateListener sent error")
		}

		err = s.processCheckSaleStateOperation(ctx, *checkSaleStateOp)
		if err != nil {
			return errors.Wrap(err, "failed to process CheckSaleState operation")
		}

		return nil
	}
}

func (s *Service) processCheckSaleStateOperation(ctx context.Context, checkSaleStateOperation horizon.CheckSaleState) error {
	if checkSaleStateOperation.Effect != xdr.CheckSaleStateEffectUpdated.String() {
		return nil
	}

	txID := checkSaleStateOperation.TransactionID
	tx, err := s.transactionConnector.TransactionByID(txID)
	if err != nil {
		return errors.Wrap(err, "failed to get transaction", logan.F{
			"transaction_id": txID,
			"operation_id":   checkSaleStateOperation.ID,
		})
	}
	if tx == nil {
		// Transaction doesn't exist
		return nil
	}

	ledgerChanges := tx.LedgerChanges()

	for _, change := range ledgerChanges {
		err = s.processCancelledOrder(ctx, change, checkSaleStateOperation.SaleID)
		if err != nil {
			return errors.Wrap(err, "failed to process cancelled order")
		}
	}

	return nil
}

func (s *Service) processCancelledOrder(ctx context.Context, change xdr.LedgerEntryChange, saleID uint64) error {
	if change.Type != xdr.LedgerEntryChangeTypeRemoved {
		return nil
	}

	removedEntryKey := change.Removed
	if removedEntryKey.Type != xdr.LedgerEntryTypeOfferEntry {
		return nil
	}
	offer := removedEntryKey.MustOffer()
	ownerID := offer.OwnerId.Address()

	user, err := s.userConnector.User(ownerID)
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

	uniqueToken := fmt.Sprintf("%s:%d:%d:%s", emailAddress, removedEntryKey.MustOffer().OfferId, saleID, s.config.RequestTokenSuffix)

	sale, err := s.saleConnector.SaleByID(saleID)
	if err != nil {
		return errors.Wrap(err, "failed to load sale", logan.F{
			"sale_id": saleID,
		})
	}
	if sale == nil {
		// Sale doesn't exist
		return nil
	}
	if sale.Name() == "" {
		return errors.New(fmt.Sprintf("invalid sale name for id: %d", saleID))
	}

	err = s.sendEmail(ctx, emailAddress, sale.Name(), uniqueToken)
	if err != nil {
		return errors.Wrap(err, "failed to send email", logan.F{
			"user_email_address": emailAddress,
		})
	}

	return nil
}
