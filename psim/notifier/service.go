package notifier

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"golang.org/x/net/context"
	"time"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/distributed_lab/logan/v3/errors"
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
	createdKYCNotifier    CreatedKYCNotifier
}

// New is a constructor of a service
func New(
	config Config,
	logger *logan.Entry,
	notificatorConnector NotificatorConnector,
	saleConnector SaleConnector,
	templatesConnector TemplatesConnector,
	transactionConnector TransactionConnector,
	userConnector UserConnector,
	checkSaleStateResponses <-chan horizon.CheckSaleStateResponse,
	createKYCRequestOpResponses <-chan horizon.CreateKYCRequestOpResponse,
) (*Service, error) {

	cancelledSaleEmailSender, err := NewOpEmailSender(
		config.SaleCancelled.Emails.Subject,
		config.SaleCancelled.Emails.TemplateName,
		config.SaleCancelled.Emails.RequestType,
		logger,
		notificatorConnector,
		templatesConnector,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cancelledSaleEmailSender")
	}

	createdKYCEmailSender, err := NewOpEmailSender(
		config.KYCCreated.Emails.Subject,
		config.KYCCreated.Emails.TemplateName,
		config.KYCCreated.Emails.RequestType,
		logger,
		notificatorConnector,
		templatesConnector,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create createdKYCEmailSender")
	}

	return &Service{
		config: config,
		logger: logger,

		cancelledSaleNotifier: CancelledSaleNotifier{
			emailSender:             cancelledSaleEmailSender,
			eventConfig:             config.SaleCancelled,
			saleConnector:           saleConnector,
			transactionConnector:    transactionConnector,
			userConnector:           userConnector,
			checkSaleStateResponses: checkSaleStateResponses,
		},

		createdKYCNotifier: CreatedKYCNotifier{
			emailSender:                 createdKYCEmailSender,
			eventConfig:                 config.KYCCreated,
			transactionConnector:        transactionConnector,
			userConnector:               userConnector,
			createKYCRequestOpResponses: createKYCRequestOpResponses,
		},
	}, nil
}

func (s *Service) Run(ctx context.Context) {
	s.logger.Info("Starting")
	//running.WithBackOff(ctx, s.logger, "cancelled_sale_notifier",
	//	s.cancelledSaleNotifier.listenAndProcessCancelledSales, 0, 5*time.Second, time.Second)
	running.WithBackOff(ctx, s.logger, "created_kyc_notifier",
		s.createdKYCNotifier.listenAndProcessCreatedKYCRequests, 0, 5*time.Second, time.Second)
}
