package supervisor

import (
	"net"

	"context"

	"time"

	"sync"

	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/psim/app"
)

// Service is common Supervisor for using in different specific Supervisors.
type Service struct {
	Log     *logan.Entry
	Errors  chan error
	Horizon *horizon.Connector

	IsLeader bool

	config    Config
	discovery *discovery.Client
	listener  net.Listener
	runners   []func(context.Context)
	builder   *xdrbuild.Builder
}

// InitNew prepares new Service (Supervisor), initializing it with all necessary helpers, got from ctx.
func InitNew(ctx context.Context, serviceName string, config Config) (*Service, error) {
	log := app.Log(ctx).WithField("service", serviceName)

	globalConfig := app.Config(ctx)

	horizonConnector := globalConfig.Horizon().WithSigner(config.SignerKP)

	listener, err := ape.Listener(config.Host, config.Port)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init listener")
	}

	// init transaction builder
	var builder *xdrbuild.Builder
	{
		info, err := horizonConnector.Info()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get horizon info")
		}
		builder = xdrbuild.NewBuilder(info.Passphrase, info.TXExpirationPeriod)
	}

	result := New(log, horizonConnector, globalConfig.Discovery(), config, listener, builder)

	result.initCommonRunners()
	return result, nil
}

func New(
	log *logan.Entry, horizon *horizon.Connector, discovery *discovery.Client, config Config, listener net.Listener,
	builder *xdrbuild.Builder,
) *Service {

	return &Service{
		Log:     log,
		Horizon: horizon,

		discovery: discovery,
		config:    config,
		listener:  listener,
		builder:   builder,
	}
}

// FIXME Discovery is switched off now.
func (s *Service) initCommonRunners() {
	if s.config.Pprof {
		s.AddRunner(s.debugAPI)
	}

	// FIXME
	// FIXME
	// FIXME
	// FIXME
	// FIXME
	// FIXME
	// FIXME
	// FIXME
	// FIXME
	// FIXME
	// FIXME
	// FIXME
	// FIXME
	// FIXME
	//s.AddRunner(s.acquireLeadership)
}

// AddRunner adds a runner to be run in separate goroutine each.
// Runner must be blocking, once runner returned - it won't be called again.
// TODO runner func must receive ctx
func (s *Service) AddRunner(runner func(context.Context)) {
	s.runners = append(s.runners, runner)
}

// Run starts all runners in separate goroutines and creates routine, which waits for all of the runners to return.
// Once all runners returned - method will finish.
// Implements app.Service.
func (s *Service) Run(ctx context.Context) {
	s.Log.Info("Started.")

	wg := sync.WaitGroup{}

	for _, runner := range s.runners {
		ohigo := runner
		wg.Add(1)

		go func() {
			ohigo(ctx)
			wg.Done()
		}()
	}

	wg.Wait()
}

func (s *Service) debugAPI(ctx context.Context) {
	s.Log.Info("enabling debug endpoints")

	r := ape.DefaultRouter()
	ape.InjectPprof(r)
	s.Log.WithField("address", s.listener.Addr().String()).Info("Starting debug API listening.")

	err := ape.ListenAndServe(ctx, s.listener, r)
	if err != nil {
		s.Log.WithError(err).Error("ListenAndServe of debug API has been stopped.")
		return
	}

	// Yes, ape.ListenAndServe can return nil error (in case of successful shutdown).
	return
}

// TODO Return error
// TODO Run over incremental timer.
func (s *Service) acquireLeadership(ctx context.Context) {
	var session *discovery.Session
	var err error

	// FIXME Select from ticker and ctx.Done() simultaneously
	ticker := time.NewTicker(5 * time.Second)
	for ; true; <-ticker.C {
		if app.IsCanceled(ctx) {
			return
		}

		if session == nil {
			session, err = discovery.NewSession(s.discovery)
			if err != nil {
				s.Log.WithError(err).Error("Failed to register session in Discovery")
				continue
			}
			session.EndlessRenew()
		}

		ok, err := s.discovery.TryAcquire(&discovery.KVPair{
			Key:     s.config.LeadershipKey,
			Session: session,
		})

		if err != nil {
			s.Log.WithError(err).Error("Failed to acquire leadership")
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
