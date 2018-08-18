package notifier

import (
	"sync"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/regources"
	"golang.org/x/net/context"
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
	TransactionByID(txID string) (*regources.Transaction, error)
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
	getKYCFirstName(kycData map[string]interface{}) (string, error)
}

type Service struct {
	config Config
	logger *logan.Entry

	cancelledOrderNotifier     CancelledOrderNotifier
	createdKYCNotifier         CreatedKYCNotifier
	reviewedKYCRequestNotifier ReviewedKYCRequestNotifier
	paymentV2Notifier          PaymentV2Notifier
}

// New is a constructor of a service
func New(
	config Config,
	log *logan.Entry,
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
	paymentV2OpResponses <-chan horizon.PaymentV2OpResponse,
) (*Service, error) {

	cancelledOrderEmailSender, err := NewOpEmailSender(
		config.OrderCancelled.Emails.Subject,
		config.OrderCancelled.Emails.TemplateName,
		config.OrderCancelled.Emails.RequestType,
		log.WithField("emails_type", "cancelled_order"),
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
		log.WithField("emails_type", "created_kyc"),
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
		log.WithField("emails_type", "approved_kyc"),
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
		log.WithField("emails_type", "rejected_kyc"),
		notificatorConnector,
		templatesConnector,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create rejectedKYCEmailSender")
	}

	usaKYCEmailSender, err := NewOpEmailSender(
		config.USAKyc.Emails.Subject,
		config.USAKyc.Emails.TemplateName,
		config.USAKyc.Emails.RequestType,
		log.WithField("emails_type", "usa_kyc"),
		notificatorConnector,
		templatesConnector,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create usaKYCEmailSender")
	}

	paymentV2EmailSender, err := NewOpEmailSender(
		config.PaymentV2.Emails.Subject,
		config.PaymentV2.Emails.TemplateName,
		config.PaymentV2.Emails.RequestType,
		log.WithField("emails_type", "payment_v2"),
		notificatorConnector,
		templatesConnector,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create paymentV2EmailSender")
	}

	return &Service{
		config: config,
		logger: log,

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
			log: log,
			approvedKYCEmailSender:   approvedKYCEmailSender,
			rejectedKYCEmailSender:   rejectedKYCEmailSender,
			usaKYCEmailSender:        usaKYCEmailSender,
			approvedRequestConfig:    config.KYCApproved,
			usaKYCConfig:             config.USAKyc,
			rejectedRequestConfig:    config.KYCRejected,
			requestConnector:         requestConnector,
			userConnector:            userConnector,
			kycDataHelper:            &KYCDataGetter{blobsConnector: blobsConnector},
			reviewRequestOpResponses: reviewRequestOpResponses,
		},

		paymentV2Notifier: PaymentV2Notifier{
			log:                  log,
			emailSender:          paymentV2EmailSender,
			eventConfig:          config.PaymentV2,
			transactionConnector: transactionConnector,
			userConnector:        userConnector,
			paymentV2Responses:   paymentV2OpResponses,
		},
	}, nil
}

func (s *Service) Run(ctx context.Context) {
	s.logger.Info("Starting")

	var opNotifiersWaitGroup sync.WaitGroup

	opNotifiersWaitGroup.Add(4)

	go func(w *sync.WaitGroup) {
		defer w.Done()
		running.WithBackOff(ctx, s.logger, "cancelled_order_notifier",
			s.cancelledOrderNotifier.listenAndProcessCancelledOrders, 0, 5*time.Second, time.Second)
	}(&opNotifiersWaitGroup)

	go func(w *sync.WaitGroup) {
		defer w.Done()
		running.WithBackOff(ctx, s.logger, "created_kyc_notifier",
			s.createdKYCNotifier.listenAndProcessCreatedKYCRequests, 0, 5*time.Second, time.Second)
	}(&opNotifiersWaitGroup)

	go func(w *sync.WaitGroup) {
		defer w.Done()
		running.WithBackOff(ctx, s.logger, "reviewed_kyc_notifier",
			s.reviewedKYCRequestNotifier.listenAndProcessReviewedKYCRequests, 0, 5*time.Second, time.Second)
	}(&opNotifiersWaitGroup)

	go func(w *sync.WaitGroup) {
		defer w.Done()
		running.WithBackOff(ctx, s.logger, "payment_v2_notifier",
			s.paymentV2Notifier.listenAndProcessPaymentV2, 0, 5*time.Second, time.Second)
	}(&opNotifiersWaitGroup)

	opNotifiersWaitGroup.Wait()
}
