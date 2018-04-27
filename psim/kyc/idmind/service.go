package idmind

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/kyc"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/keypair"
)

// RequestListener is the interface, which must be implemented
// by streamer of Horizon Requests, which parametrize Service.
type RequestListener interface {
	StreamAllKYCRequests(ctx context.Context, endlessly bool) <-chan horizon.ReviewableRequestEvent
	StreamKYCRequestsUpdatedAfter(ctx context.Context, updatedAfter time.Time, endlessly bool) <-chan horizon.ReviewableRequestEvent
}

type RequestPerformer interface {
	Approve(ctx context.Context, requestID uint64, requestHash string, tasksToAdd, tasksToRemove uint32, extDetails map[string]string) error
	Reject(ctx context.Context, requestID uint64, requestHash string, tasksToAdd uint32, extDetails map[string]string, rejectReason string) error
}

type BlobSubmitter interface {
	SubmitBlob(ctx context.Context, blobType, attrValue string, relationships map[string]string) (blobID string, err error)
}

type BlobDataRetriever interface {
	RetrieveKYCBlob(kycRequest horizon.KYCRequest, accountID string) (*horizon.Blob, error)
}

type DocumentsConnector interface {
	Document(docID string) (*horizon.Document, error)
}

type UsersConnector interface {
	User(accountID string) (*horizon.User, error)
}

type IdentityMind interface {
	Submit(req CreateAccountRequest) (*ApplicationResponse, error)
	CheckState(txID string) (*CheckApplicationResponse, error)
}

type EmailsProcessor interface {
	Run(ctx context.Context)
	AddEmailAddresses(ctx context.Context, subject, message string, emailAddresses []string)
	AddTask(ctx context.Context, emailAddress, subject, message string)
}

type Service struct {
	log    *logan.Entry
	config Config
	signer keypair.Full
	source keypair.Address

	requestListener    RequestListener
	requestPerformer   RequestPerformer
	blobsConnector     BlobSubmitter
	blobDataRetriever  BlobDataRetriever
	documentsConnector DocumentsConnector
	usersConnector     UsersConnector
	identityMind       IdentityMind
	adminNotifyEmails  EmailsProcessor

	kycRequests <-chan horizon.ReviewableRequestEvent
}

// NewService is constructor for Service.
func NewService(
	log *logan.Entry,
	config Config,
	requestListener RequestListener,
	requestPerformer RequestPerformer,
	blobProvider BlobSubmitter,
	blobDataRetriever BlobDataRetriever,
	userProvider UsersConnector,
	documentProvider DocumentsConnector,
	identityMind IdentityMind,
	adminNotifyEmails EmailsProcessor,
) *Service {

	return &Service{
		log:    log.WithField("service", conf.ServiceIdentityMind),
		config: config,

		requestListener:    requestListener,
		requestPerformer:   requestPerformer,
		blobsConnector:     blobProvider,
		blobDataRetriever:  blobDataRetriever,
		usersConnector:     userProvider,
		documentsConnector: documentProvider,
		identityMind:       identityMind,
		adminNotifyEmails:  adminNotifyEmails,
	}
}

// Run is a blocking method, it returns only when ctx closes.
func (s *Service) Run(ctx context.Context) {
	s.log.WithField("", s.config).Info("Starting.")

	go s.adminNotifyEmails.Run(ctx)
	s.kycRequests = s.requestListener.StreamAllKYCRequests(ctx, false)

	running.WithBackOff(ctx, s.log, "kyc_request_processor", s.listenAndProcessRequest, 0, 5*time.Second, 5*time.Minute)
}

// TODO timeToSleep to config
func (s *Service) listenAndProcessRequest(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	case reqEvent, ok := <-s.kycRequests:
		if !ok {
			// No more KYC requests, start from the very beginning.
			// TODO timeToSleep to config
			timeToSleep := 30 * time.Second
			s.log.Debugf("No more KYC Requests in Horizon, will start from the very beginning, now sleeping for (%s).", timeToSleep.String())

			c := time.After(timeToSleep)
			select {
			case <-ctx.Done():
				return nil
			case <-c:
				s.kycRequests = s.requestListener.StreamAllKYCRequests(ctx, false)
				return nil
			}
		}

		request, err := reqEvent.Unwrap()
		if err != nil {
			return errors.Wrap(err, "RequestListener sent error")
		}

		err = s.processRequest(ctx, *request)
		if err != nil {
			return errors.Wrap(err, "Failed to process KYC Request", logan.F{
				"request": request,
			})
		}

		return nil
	}
}

func (s *Service) processRequest(ctx context.Context, request horizon.Request) error {
	proveErr := proveInterestingRequest(request)
	if proveErr != nil {
		// No need to process the Request for now.

		// I found this log useless
		//s.log.WithField("request", request).WithError(proveErr).Debug("Found not interesting KYC Request.")
		return nil
	}

	// I found this log useless
	s.log.WithField("request", request).Debug("Found interesting KYC Request.")
	kycReq := request.Details.KYC

	if kycReq.PendingTasks&kyc.TaskSubmitIDMind != 0 {
		// Haven't submitted IDMind yet
		err := s.processNotSubmitted(ctx, request)
		if err != nil {
			return errors.Wrap(err, "Failed to process not submitted (to IDMind) KYCRequest")
		}

		return nil
	}

	// Already submitted
	if kycReq.PendingTasks&kyc.TaskCheckIDMind != 0 {
		err := s.processNotChecked(ctx, request)
		if err != nil {
			return errors.Wrap(err, "Failed to check KYC state in IDMind")
		}

		return nil
	}

	// Normally unreachable
	return nil
}
