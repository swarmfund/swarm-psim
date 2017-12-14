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

// Service is common Supervisor for using in different specific Supervisors.
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

// InitNew prepares new Service (Supervisor), initializing it with all necessary helpers, got from ctx.
func InitNew(ctx context.Context, serviceName string, config Config) (*Service, error) {
	log := app.Log(ctx).WithField("service", serviceName)

	globalConfig := app.Config(ctx)

	discoveryClient, err := globalConfig.Discovery()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get DiscoveryClient")
	}

	horizonConnector, err := globalConfig.Horizon()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get HorizonClient")
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
		// TODO Avoid ctx from in struct
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
// TODO runner func must receive ctx
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
				// TODO runner func must receive ctx
				ohigo()
				wg.Done()
			}()
		}

		wg.Wait()
		close(s.Errors)
	}()

	return s.Errors
}

// TODO runner func must receive ctx
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

// TODO runner func must receive ctx
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
