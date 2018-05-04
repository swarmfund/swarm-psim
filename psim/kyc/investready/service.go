package investready

import (
	"context"
	"time"

	"sync"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/keypair"
	"gitlab.com/swarmfund/psim/psim/kyc"
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

type BlobsConnector interface {
	Blob(blobID string) (*horizon.Blob, error)
	SubmitBlob(ctx context.Context, blobType, attrValue string, relationships map[string]string) (blobID string, err error)
}

type BlobDataRetriever interface {
	ParseBlobData(kycRequest horizon.KYCRequest) (*kyc.Data, error)
}

type UsersConnector interface {
	User(accountID string) (*horizon.User, error)
}

type KYCRequestsConnector interface {
	Requests(filters, cursor string, reqType horizon.ReviewableRequestType) ([]horizon.Request, error)
}

type InvestReady interface {
	ObtainUserToken(oauthCode string) (userAccessToken string, err error)
	UserHash(userAccessToken string) (userHash string, err error)
	ListAllSyncedUsers(ctx context.Context) ([]User, error)
}

// TODO Comment
type Service struct {
	log    *logan.Entry
	config Config
	signer keypair.Full
	source keypair.Address

	requestListener  RequestListener
	requestPerformer RequestPerformer
	blobsConnector   BlobsConnector
	blobDataRetriever BlobDataRetriever
	usersConnector   UsersConnector

	investReady InvestReady

	redirectsListener *RedirectsListener
	kycRequests       <-chan horizon.ReviewableRequestEvent
	users             []User
}

// TODO add docs.md

// NewService is constructor for Service.
func NewService(
	log *logan.Entry,
	config Config,
	requestListener RequestListener,
	kycRequestsConnector KYCRequestsConnector,
	requestPerformer RequestPerformer,
	blobProvider BlobsConnector,
	blobDataRetriever BlobDataRetriever,
	userProvider UsersConnector,
	investReady InvestReady,
) *Service {

	logger := log.WithField("service", conf.ServiceInvestReady)
	return &Service{
		log:    logger,
		config: config,

		requestListener:  requestListener,
		requestPerformer: requestPerformer,
		blobsConnector:   blobProvider,
		blobDataRetriever: blobDataRetriever,
		usersConnector:   userProvider,
		investReady:      investReady,

		redirectsListener: NewRedirectsListener(logger, config.RedirectsConfig, kycRequestsConnector, investReady, requestPerformer),
	}
}

// Run is a blocking method, it returns only when ctx closes.
func (s *Service) Run(ctx context.Context) {
	s.log.WithField("", s.config).Info("Starting.")

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		s.redirectsListener.Run(ctx)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		s.processRequestsInfinitely(ctx)
		wg.Done()
	}()

	wg.Wait()
	s.log.Info("All runners stopped - stopping cleanly.")
}
