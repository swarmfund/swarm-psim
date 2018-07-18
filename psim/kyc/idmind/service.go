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

var (
	errInvalidState = errors.New("service got into unknown state")
)

// RequestListener is the interface, which must be implemented
// by streamer of Horizon Requests, which parametrize Service.
type RequestListener interface {
	StreamAllKYCRequests(ctx context.Context, endlessly bool) <-chan horizon.ReviewableRequestEvent
	StreamKYCRequestsUpdatedAfter(ctx context.Context, updatedAfter time.Time, endlessly bool) <-chan horizon.ReviewableRequestEvent
}

type RequestPerformer interface {
	Approve(ctx context.Context, requestID uint64, requestHash string, tasksToAdd, tasksToRemove uint32, extDetails map[string]string) error
	Reject(ctx context.Context, requestID uint64, requestHash string, tasksToAdd uint32, extDetails map[string]string, rejectReason, rejector string) error
}

type BlobSubmitter interface {
	SubmitBlob(ctx context.Context, blobType, attrValue string, relationships map[string]string) (blobID string, err error)
}

type DocumentsConnector interface {
	Document(docID string) (*horizon.Document, error)
}

type UsersConnector interface {
	User(accountID string) (*horizon.User, error)
}

type AccountsConnector interface {
	ByAddress(address string) (*horizon.Account, error)
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

	horizon            *horizon.Connector
	requestListener    RequestListener
	requestPerformer   RequestPerformer
	blobSubmitter      BlobSubmitter
	documentsConnector DocumentsConnector
	usersConnector     UsersConnector
	accountsConnector  AccountsConnector
	identityMind       IdentityMind
	adminNotifyEmails  EmailsProcessor

	kycRequests <-chan horizon.ReviewableRequestEvent
}

// NewService is constructor for Service.
func NewService(
	log *logan.Entry,
	config Config,
	horizon *horizon.Connector,
	requestListener RequestListener,
	requestPerformer RequestPerformer,
	blobSubmitter BlobSubmitter,
	usersConnector UsersConnector,
	accountsConnector AccountsConnector,
	documentProvider DocumentsConnector,
	identityMind IdentityMind,
	adminNotifyEmails EmailsProcessor,
) *Service {

	return &Service{
		log:    log.WithField("service", conf.ServiceIdentityMind),
		config: config,

		horizon:            horizon,
		requestListener:    requestListener,
		requestPerformer:   requestPerformer,
		blobSubmitter:      blobSubmitter,
		usersConnector:     usersConnector,
		accountsConnector:  accountsConnector,
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
	fields := logan.F{
		"request": request,
	}

	// check if request should be processed
	if ok := isInterestingRequest(request); !ok {
		s.log.WithFields(fields).Debug("skipping not interesting request")
		return nil
	}

	// check if request is valid
	if err := isValidRequest(request); err != nil {
		return errors.Wrap(err, "request is invalid in some way")
	}

	kycDetails := request.Details.KYC

	// check if account we are going to review is blocked
	account, err := s.accountsConnector.ByAddress(kycDetails.AccountToUpdateKYC)
	if err != nil {
		return errors.Wrap(err, "failed to get account", logan.F{
			"address": kycDetails.AccountToUpdateKYC,
		})
	}
	if account.IsBlocked {
		s.log.WithFields(fields).Debug("skipping since account is blocked")
		return nil
	}

	// check if submitted blob is valid
	blob, err := s.horizon.Blobs().Blob(kycDetails.KYCDataStruct.BlobID)
	if err != nil {
		return errors.Wrap(err, "failed to get blob", logan.F{
			"blob_id": kycDetails.KYCDataStruct.BlobID,
		})
	}
	if err := isBlobValid(blob); err != nil {
		return errors.Wrap(err, "blob is invalid in some way")
	}
	kycData, err := kyc.ParseKYCData(blob.Attributes.Value)
	if err != nil {
		return errors.Wrap(err, "failed to parse KYC data", logan.F{
			"blob_id": blob.ID,
		})
	}

	// check if we are able process given account type transition
	switch {
	case kyc.IsNotVerified(*account) && kyc.IsUpdateToGeneral(request) && !kycData.IsUSA():
	case kyc.IsNotVerified(*account) && kyc.IsUpdateToVerified(request) && kycData.IsUSA():
	case kyc.IsGeneral(*account) && kyc.IsUpdateToVerified(request) && kycData.IsUSA():
	case kyc.IsGeneral(*account) && kyc.IsUpdateToGeneral(request) && !kycData.IsUSA():
	default:
		err := s.rejectRequest(ctx, request, s.config.RejectReasons.KYCStateRejected, nil)
		if err != nil {
			return errors.Wrap(err, "failed to reject request")
		}
	}

	s.log.WithFields(fields).Debug("processing KYC request")

	// finally process request
	switch {
	case kycDetails.PendingTasks&kyc.TaskNonLatinDoc != 0:
		if err := s.approveRequest(ctx, request, nil); err != nil {
			return errors.Wrap(err, "failed to approve request with a non-latin-docs task set")
		}
	case kycDetails.PendingTasks&kyc.TaskSubmitIDMind != 0:
		err = s.processNewKYCApplication(ctx, kycData, request)
	case kycDetails.PendingTasks&kyc.TaskCheckIDMind != 0:
		err = s.processNotChecked(ctx, request)
	default:
		// unknown state, probably someone messed up isInterestingRequest call above
		return errInvalidState
	}
	if err != nil {
		return errors.Wrap(err, "failed to process IDMind request", logan.F{
			"pending_tasks": kycDetails.PendingTasks,
		})
	}

	s.log.WithFields(fields).Debug("Processed KYC request successfully.")

	return nil
}
