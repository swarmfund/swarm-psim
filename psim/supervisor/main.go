package supervisor

import (
	"net"

	"context"

	"time"

	"sync"

	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/psim/app"
)

// TODO Add comment
type Service struct {
	Ctx    context.Context
	Log    *logan.Entry
	Errors chan error

	IsLeader bool

	config Config
	// TODO interface?
	horizon   *horizon.Connector
	discovery *discovery.Client
	listener  net.Listener
	runners   []func()
}

// TODO Add comment
func InitNew(ctx context.Context, serviceName string, config Config) (*Service, error) {
	log := app.Log(ctx).WithField("service", serviceName)

	globalConfig := app.Config(ctx)

	discoveryClient, err := globalConfig.Discovery()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get discovery client")
	}

	horizonConnector, err := globalConfig.Horizon()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get horizon client")
	}

	listener, err := ape.Listener(config.Host, config.Port)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init listener")
	}

	result := New(ctx, log, horizonConnector, discoveryClient, config, listener)

	result.initCommonRunners()
	return result, nil
}

func New(ctx context.Context, log *logan.Entry, horizon *horizon.Connector, discovery *discovery.Client, config Config, listener net.Listener) *Service {

	return &Service{
		Ctx:    ctx,
		Log:    log,
		Errors: make(chan error),

		horizon:   horizon,
		discovery: discovery,
		config:    config,
		listener:  listener,
	}
}

func (s *Service) initCommonRunners() {
	if s.config.Pprof {
		s.AddRunner(s.debugAPI)
	}
	s.AddRunner(s.acquireLeadership)
}

// AddRunner adds a runner to be run in separate goroutine each.
// Runner must be blocking, once runner returned - it won't be called again.
func (s *Service) AddRunner(runner func()) {
	s.runners = append(s.runners, runner)
}

// Run starts all runners in separate goroutines and creates routine, which waits for all of the runners to return.
// Once all runners returned - Errors channel will be closed.
// Implements utils.Service.
func (s *Service) Run() chan error {
	go func() {
		wg := sync.WaitGroup{}

		for _, runner := range s.runners {
			ohigo := runner
			wg.Add(1)

			go func() {
				ohigo()
				wg.Done()
			}()
		}

		wg.Wait()
		close(s.Errors)
	}()

	return s.Errors
}

func (s *Service) debugAPI() {
	s.Log.Info("enabling debug endpoints")

	r := ape.DefaultRouter()
	ape.InjectPprof(r)
	s.Log.WithField("address", s.listener.Addr().String()).Info("listening")

	err := ape.ListenAndServe(s.Ctx, s.listener, r)
	if err != nil {
		s.Errors <- err
		return
	}
	return
}

func (s *Service) acquireLeadership() {
	var session *discovery.Session
	var err error

	// FIXME Select from ticker and ctx.Done() simultaneously
	ticker := time.NewTicker(5 * time.Second)
	for ; true; <-ticker.C {
		if app.IsCanceled(s.Ctx) {
			return
		}

		if session == nil {
			session, err = discovery.NewSession(s.discovery)
			if err != nil {
				s.Errors <- errors.Wrap(err, "Failed to register session in Discovery")
				continue
			}
			session.EndlessRenew()
		}

		ok, err := s.discovery.TryAcquire(&discovery.KVPair{
			Key:     s.config.LeadershipKey,
			Session: session,
		})

		if err != nil {
			s.Errors <- err
			s.IsLeader = false
			continue
		}

		if ok {
			s.IsLeader = true
		} else {
			// probably will never happen, but just in case
			s.IsLeader = false
		}
	}
}
