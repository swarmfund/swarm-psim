package notifier

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"golang.org/x/net/context"
	"time"
	"gitlab.com/distributed_lab/running"
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

type EmailSender interface {
	SendEmail(ctx context.Context, emailAddress, emailUniqueToken string, data interface{}) error
}

type Service struct {
	config Config
	logger *logan.Entry

	cancelledSaleNotifier CancelledSaleNotifier
}

// New is a constructor of a service
func New(
	config Config,
	logger *logan.Entry,
	emailSender EmailSender,
	saleConnector SaleConnector,
	transactionConnector TransactionConnector,
	userConnector UserConnector,
	checkSaleStateResponses <-chan horizon.CheckSaleStateResponse,
) *Service {
	return &Service{
		config: config,
		logger: logger,

		cancelledSaleNotifier: CancelledSaleNotifier{
			emailSender:             emailSender,
			emailsConfig:            config.SaleCancelled,
			saleConnector:           saleConnector,
			transactionConnector:    transactionConnector,
			userConnector:           userConnector,
			checkSaleStateResponses: checkSaleStateResponses,
		},
	}
}

func (s *Service) Run(ctx context.Context) {
	s.logger.Info("Starting...")
	go running.WithBackOff(ctx, s.logger, "cancelled_sale_notifier",
		s.cancelledSaleNotifier.listenAndProcessCancelledSales, 0, 5*time.Second, time.Second)
}
