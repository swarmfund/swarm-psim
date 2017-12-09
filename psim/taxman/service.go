package taxman

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	discovery "gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	sse "gitlab.com/distributed_lab/sse-go"
	horizon "gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/psim/psim/taxman/internal/snapshoter"
	"gitlab.com/swarmfund/psim/psim/taxman/internal/state"
	"gitlab.com/swarmfund/psim/psim/taxman/internal/txhandler"
	"gitlab.com/swarmfund/psim/psim/utils"
)

type Service struct {
	ID string

	log              *logan.Entry
	discovery        *discovery.Client
	discoveryService *discovery.Service
	horizon          *horizon.Connector
	errors           chan error
	sse              *sse.Listener
	listener         net.Listener
	state            *state.State
	txHandler        *txhandler.Handler
	snapshots        snapshoter.Snapshots
	ticker           *time.Ticker
	server           *http.Server
	config           Config

	TxCursor   string
	NextPayout time.Time

	// teardown
	ctx    context.Context
	cancel context.CancelFunc
}

func New(log *logan.Entry, discovery *discovery.Client, horizon *horizon.Connector,
	listener net.Listener, config Config, ctx context.Context,
) *Service {
	payoutState := state.NewState()
	ctx, cancel := context.WithCancel(ctx)
	service := Service{
		ID:        utils.GenerateToken(),
		log:       log,
		discovery: discovery,
		horizon:   horizon,
		listener:  listener,
		state:     payoutState,
		txHandler: txhandler.NewHandler(payoutState, config.Skip, log),
		errors:    make(chan error),
		snapshots: snapshoter.Snapshots{},
		ticker:    time.NewTicker(5 * time.Second),
		config:    config,
		cancel:    cancel,
		ctx:       ctx,
	}

	service.sse = sse.NewListener(func() (*http.Request, error) {
		return horizon.SignedRequest("GET", fmt.Sprintf("/transactions?cursor=%s", service.TxCursor), config.Signer)
	})

	return &service
}

func (s *Service) TearDown() {
	s.log.Debug("tearing down")
	s.cancel()
}

func (s *Service) Run() chan error {
	wg := sync.WaitGroup{}
	serviceSeq := []func(ctx context.Context){
		//s.AcquireLeadership,
		s.API,
		s.Register,
		s.Listener,
	}

	for _, fn := range serviceSeq {
		f := fn
		wg.Add(1)
		go func() {
			defer func() {
				if rec := recover(); rec != nil {
					s.errors <- errors.FromPanic(rec)
				}
				wg.Done()
			}()
			f(s.ctx)
		}()
	}

	go func() {
		defer func() {
			close(s.errors)
		}()
		wg.Wait()
	}()

	return s.errors
}

func (s *Service) SetUp(ctx context.Context) error {
	return nil
}
