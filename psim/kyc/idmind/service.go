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
		"request": request.GetLoganFields(),
	}

	// check if request should be processed
	if ok := isInterestingRequest(request); !ok {
		s.log.WithFields(fields).Debug("skipping not interesting request")
		return nil
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
		return nil
	}

	// check if submitted blob is valid
	blob, err := s.horizon.Blobs().Blob(kycDetails.KYCDataStruct.BlobID)
	if err != nil {
		return errors.Wrap(err, "failed to get blob", logan.F{
			"blob_id": kycDetails.KYCDataStruct.BlobID,
		})
	}
	if ok := isBlobValid(blob); !ok {
		// TODO consider rejecting invalid requests
		s.log.WithFields(fields).Debug("skipping request with malformed blob")
		return nil
	}
	kycData, err := kyc.ParseKYCData(blob.Attributes.Value)
	if err != nil {
		return errors.Wrap(err, "failed to parse KYC data", logan.F{
			"blob_id": blob.ID,
		})
	}

	// check if we are able process given account type transition
	switch {
	case isNotVerified(account) && isUpdateToGeneral(request) && !kycData.IsUSA():
	case isNotVerified(account) && isUpdateToVerified(request) && kycData.IsUSA():
	case isGeneral(account) && isUpdateToVerified(request) && kycData.IsUSA():
	case isGeneral(account) && isUpdateToGeneral(request) && !kycData.IsUSA():
	default:
		// TODO reject request as not valid
		panic("not implemented")
	}

	s.log.WithFields(fields).Debug("processing KYC request")

	// finally process request
	switch {
	// FIXME (stepko) why all instead of pending?
	case kycDetails.AllTasks&kyc.TaskNonLatinDoc != 0:
		// TODO approve our tasks w/o submit
		panic("not implemented")
	case kycDetails.PendingTasks&kyc.TaskSubmitIDMind != 0:
		err = s.processNewKYCApplication(ctx, kycData, request)
	case kycDetails.PendingTasks&kyc.TaskCheckIDMind != 0:
		err = s.processNotChecked(ctx, request)
	default:
		// unknown state, probably some one messed up isInterestingRequest call above
		panic("not implemented")
	}
	if err != nil {
		return errors.Wrap(err, "failed to process IDMind request", logan.F{
			"pending_tasks": kycDetails.PendingTasks,
		})
	}

	s.log.WithFields(fields).Debug("processed KYC request successfully")

	return nil
}
