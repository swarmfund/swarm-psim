package investready

import (
	"context"
	"time"

	"sync"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/keypair"
	"gitlab.com/distributed_lab/running"
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

type BlobsConnector interface {
	Blob(blobID string) (*horizon.Blob, error)
}

type UsersConnector interface {
	User(accountID string) (*horizon.User, error)
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

	requestListener   RequestListener
	requestPerformer  RequestPerformer
	blobsConnector    BlobsConnector
	usersConnector    UsersConnector

	investReady InvestReady

	redirectsListener *RedirectsListener
	kycRequests       <-chan horizon.ReviewableRequestEvent
	syncedUserHashes  []User
}

// TODO add docs.md

// NewService is constructor for Service.
func NewService(
	log *logan.Entry,
	config Config,
	requestListener RequestListener,
	kycRequestsConnector KYCRequestsConnector,
	accountsConnector AccountsConnector,
	requestPerformer RequestPerformer,
	blobProvider BlobsConnector,
	userProvider UsersConnector,
	investReady InvestReady,
	doorman doorman.Doorman) *Service {

	logger := log.WithField("service", conf.ServiceInvestReady)
	redirectsListener := NewRedirectsListener(
		logger,
		config.RedirectsConfig,
		kycRequestsConnector,
		accountsConnector,
		investReady,
		doorman,
		requestPerformer)

	return &Service{
		log:    logger,
		config: config,

		requestListener:  requestListener,
		requestPerformer: requestPerformer,
		blobsConnector:   blobProvider,
		usersConnector:   userProvider,
		investReady:      investReady,

		redirectsListener: redirectsListener,
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

	// TODO period to config
	period := 30 * time.Second
	wg.Add(1)
	go func() {
		running.WithBackOff(ctx, s.log, "requests_processing_iteration", s.processAllRequestsOnce, period, period, period)
		wg.Done()
	}()

	wg.Wait()
	s.log.Info("All runners stopped - stopping cleanly.")
}
