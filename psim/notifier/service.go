package notifier

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"golang.org/x/net/context"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"time"
	"sync"
	"gitlab.com/swarmfund/psim/psim/kyc"
)

const (
	KYCFormBlobType = "kyc_form"
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

type KYCDataHelper interface {
	getBlobKYCData(kycData map[string]interface{}) (*kyc.Data, error)
}

type Service struct {
	config Config
	logger *logan.Entry

	cancelledOrderNotifier     CancelledOrderNotifier
	createdKYCNotifier         CreatedKYCNotifier
	reviewedKYCRequestNotifier ReviewedKYCRequestNotifier
}

// New is a constructor of a service
func New(
	config Config,
	logger *logan.Entry,
	notificatorConnector NotificatorConnector,
	requestConnector ReviewableRequestConnector,
	saleConnector SaleConnector,
	templatesConnector TemplatesConnector,
	transactionConnector TransactionConnector,
	userConnector UserConnector,
	blobsConnector BlobsConnector,
	checkSaleStateResponses <-chan horizon.CheckSaleStateResponse,
	createKYCRequestOpResponses <-chan horizon.CreateKYCRequestOpResponse,
	reviewRequestOpResponses <-chan horizon.ReviewRequestOpResponse,
) (*Service, error) {

	cancelledOrderEmailSender, err := NewOpEmailSender(
		config.OrderCancelled.Emails.Subject,
		config.OrderCancelled.Emails.TemplateName,
		config.OrderCancelled.Emails.RequestType,
		logger,
		notificatorConnector,
		templatesConnector,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cancelledOrderEmailSender")
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

	approvedKYCEmailSender, err := NewOpEmailSender(
		config.KYCApproved.Emails.Subject,
		config.KYCApproved.Emails.TemplateName,
		config.KYCApproved.Emails.RequestType,
		logger,
		notificatorConnector,
		templatesConnector,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create approvedKYCEmailSender")
	}

	rejectedKYCEmailSender, err := NewOpEmailSender(
		config.KYCRejected.Emails.Subject,
		config.KYCRejected.Emails.TemplateName,
		config.KYCRejected.Emails.RequestType,
		logger,
		notificatorConnector,
		templatesConnector,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create rejectedKYCEmailSender")
	}

	return &Service{
		config: config,
		logger: logger,

		cancelledOrderNotifier: CancelledOrderNotifier{
			emailSender:             cancelledOrderEmailSender,
			eventConfig:             config.OrderCancelled,
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
			kycDataHelper:               &KYCDataGetter{blobsConnector: blobsConnector},
			createKYCRequestOpResponses: createKYCRequestOpResponses,
		},

		reviewedKYCRequestNotifier: ReviewedKYCRequestNotifier{
			approvedKYCEmailSender:   approvedKYCEmailSender,
			rejectedKYCEmailSender:   rejectedKYCEmailSender,
			approvedRequestConfig:    config.KYCApproved,
			rejectedRequestConfig:    config.KYCRejected,
			requestConnector:         requestConnector,
			userConnector:            userConnector,
			kycDataHelper:            &KYCDataGetter{blobsConnector: blobsConnector},
			reviewRequestOpResponses: reviewRequestOpResponses,
		},
	}, nil
}

func (s *Service) Run(ctx context.Context) {
	s.logger.Info("Starting")

	var opNotifiersWaitGroup sync.WaitGroup

	opNotifiersWaitGroup.Add(3)

	go func(w *sync.WaitGroup) {
		running.WithBackOff(ctx, s.logger, "cancelled_order_notifier",
			s.cancelledOrderNotifier.listenAndProcessCancelledOrders, 0, 5*time.Second, time.Second)
		w.Done()
	}(&opNotifiersWaitGroup)

	go func(w *sync.WaitGroup) {
		running.WithBackOff(ctx, s.logger, "created_kyc_notifier",
			s.createdKYCNotifier.listenAndProcessCreatedKYCRequests, 0, 5*time.Second, time.Second)
		w.Done()
	}(&opNotifiersWaitGroup)

	go func(w *sync.WaitGroup) {
		running.WithBackOff(ctx, s.logger, "reviewed_kyc_notifier",
			s.reviewedKYCRequestNotifier.listenAndProcessReviewedKYCRequests, 0, 5*time.Second, time.Second)
		w.Done()
	}(&opNotifiersWaitGroup)

	opNotifiersWaitGroup.Wait()
}
