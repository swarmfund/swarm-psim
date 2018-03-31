package idmind

import (
	"context"
	"time"

	"io"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/tokend/keypair"
)

const (
	KYCFormBlobType = "kyc_form"
)

// RequestListener is the interface, which must be implemented
// by streamer of Horizon Requests, which parametrize Service.
type RequestListener interface {
	StreamAllKYCRequests(ctx context.Context, endlessly bool) <-chan horizon.ReviewableRequestEvent
	StreamKYCRequestsUpdatedAfter(ctx context.Context, updatedAfter time.Time, endlessly bool) <-chan horizon.ReviewableRequestEvent
}

type TXSubmitter interface {
	Submit(ctx context.Context, envelope string) horizon.SubmitResult
}

type BlobProvider interface {
	Blob(blobID string) (*horizon.Blob, error)
}

type DocumentProvider interface {
	Document(docID string) (*horizon.Document, error)
}

type UserProvider interface {
	User(accountID string) (*horizon.User, error)
}

type IdentityMind interface {
	Submit(req CreateAccountRequest) (*ApplicationResponse, error)
	UploadDocument(txID, description string, fileName string, fileReader io.Reader) error
	CheckState(txID string) (*CheckApplicationResponse, error)
}

type Service struct {
	log    *logan.Entry
	config Config
	signer keypair.Full
	source keypair.Address

	requestListener  RequestListener
	txSubmitter      TXSubmitter
	blobProvider     BlobProvider
	documentProvider DocumentProvider
	userProvider     UserProvider
	identityMind     IdentityMind
	xdrbuilder       *xdrbuild.Builder

	kycRequests <-chan horizon.ReviewableRequestEvent
}

// NewService is constructor for Service.
func NewService(
	log *logan.Entry,
	config Config,
	requestListener RequestListener,
	txSubmitter TXSubmitter,
	blobProvider BlobProvider,
	userProvider UserProvider,
	documentProvider DocumentProvider,
	identityMind IdentityMind,
	builder *xdrbuild.Builder,
) *Service {

	return &Service{
		log:    log.WithField("service", conf.ServiceIdentityMind),
		config: config,

		requestListener:  requestListener,
		txSubmitter:      txSubmitter,
		blobProvider:     blobProvider,
		userProvider:     userProvider,
		documentProvider: documentProvider,
		identityMind:     identityMind,
		xdrbuilder:       builder,
	}
}

// Run is a blocking method, it returns only when ctx closes.
func (s *Service) Run(ctx context.Context) {
	s.log.WithField("", s.config).Info("Starting.")

	s.kycRequests = s.requestListener.StreamAllKYCRequests(ctx, false)

	running.WithBackOff(ctx, s.log, "request_processor", s.listenAndProcessRequest, 0, 5*time.Second, 5*time.Minute)
}

func (s *Service) listenAndProcessRequest(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	case reqEvent, ok := <-s.kycRequests:
		if !ok {
			// No more KYC requests, start from the very beginning.
			// TODO Consider timeToSleep to config?
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
	//s.log.WithField("request", request).Debug("Found interesting KYC Request.")
	kyc := request.Details.KYC

	if kyc.PendingTasks&TaskSubmitIDMind != 0 {
		// Haven't submitted IDMind yet
		err := s.submitKYCData(ctx, request)
		if err != nil {
			return errors.Wrap(err, "Failed to submit KYC data")
		}

		return nil
	}

	if kyc.PendingTasks&TaskCheckIDMind != 0 {
		err := s.checkKYCState(ctx, request)
		if err != nil {
			return errors.Wrap(err, "Failed to check KYC state to IDMind")
		}

		return nil
	}

	// Normally unreachable
	return nil
}
