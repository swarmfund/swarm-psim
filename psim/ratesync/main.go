package ratesync

import (
	"context"
	"fmt"
	"net"
	"time"

	"reflect"

	"github.com/mitchellh/mapstructure"
	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	horizon "gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/figure"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/ratesync/noor"
	"gitlab.com/swarmfund/psim/psim/ratesync/providers"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.Register(conf.ServiceRateSync, func(ctx context.Context) error {
		serviceConfig := Config{
			ServiceName:   conf.ServiceRateSync,
			LeadershipKey: fmt.Sprintf("service/%s/leader", conf.ServiceRateSync),
			Host:          "localhost",
		}

		assetsHook := figure.Hooks{
			"ratesync.NoorConfig": func(raw interface{}) (reflect.Value, error) {
				result := NoorConfig{}
				err := mapstructure.Decode(raw, &result)
				if err != nil {
					return reflect.Value{}, err
				}
				return reflect.ValueOf(result), nil
			},
			"[]ratesync.Asset": func(raw interface{}) (reflect.Value, error) {
				result := []Asset{}
				err := mapstructure.Decode(raw, &result)
				if err != nil {
					return reflect.Value{}, err
				}
				return reflect.ValueOf(result), nil
			},
		}

		globalConfig := ctx.Value(app.CtxConfig).(conf.Config)
		err := figure.
			Out(&serviceConfig).
			From(globalConfig.Get(conf.ServiceRateSync)).
			With(figure.BaseHooks, utils.CommonHooks, assetsHook).
			Please()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to figure out %s", conf.ServiceRateSync))
		}

		retryTicker := time.NewTicker(5 * time.Second)

		log := ctx.Value(app.CtxLog).(*logan.Entry)

		listener, err := ape.Listener(serviceConfig.Host, serviceConfig.Port)
		if err != nil {
			return errors.Wrap(err, "failed to init listener")
		}

		horizonC, err := globalConfig.Horizon()
		if err != nil {
			return errors.Wrap(err, "failed to get horizon client")
		}

		providers := providers.Providers{
			"noor": noor.NewProvider(log, serviceConfig.Noor.Host, serviceConfig.Noor.Port, serviceConfig.Noor.Pairs),
		}

		provider, ok := providers[serviceConfig.Provider]
		if !ok {
			return fmt.Errorf("provider %s is not configured", serviceConfig.Provider)
		}

		service, errs := New(log, globalConfig.Discovery(), listener, horizonC, serviceConfig, provider)
		for service.Run(); ; <-retryTicker.C {
			err := <-errs
			if serr, ok := err.(horizon.SubmitError); ok {
				log.WithError(err).
					WithField("tx code", serr.TransactionCode()).
					WithField("op codes", serr.OperationCodes()).
					Error("tx failed")
			} else {
				log.WithError(<-errs).Warn("service error")
			}
		}
	})
}

type Service struct {
	ID       string
	IsLeader bool

	log              *logan.Entry
	discovery        *discovery.Client
	discoveryService *discovery.Service
	listener         net.Listener
	errors           chan error
	syncResults      *SyncResults
	syncTimestamp    chan time.Time
	horizon          *horizon.Connector
	config           Config
	provider         providers.Provider
}

func New(
	log *logan.Entry, discovery *discovery.Client, listener net.Listener, horizon *horizon.Connector,
	config Config, provider providers.Provider,
) (*Service, chan error) {
	errs := make(chan error)
	return &Service{
		ID:          utils.GenerateToken(),
		log:         log,
		discovery:   discovery,
		listener:    listener,
		errors:      errs,
		horizon:     horizon,
		config:      config,
		provider:    provider,
		syncResults: NewSyncResults(),
	}, errs
}

func (s *Service) Register() {
	s.discoveryService = s.discovery.Service(&discovery.ServiceRegistration{
		Name: s.config.ServiceName,
		ID:   s.ID,
		Host: fmt.Sprintf("http://%s/", s.listener.Addr().String()),
	})
	ticker := time.NewTicker(5 * time.Second)
	for ; true; <-ticker.C {
		err := s.discovery.RegisterServiceSync(s.discoveryService)
		if err != nil {
			s.errors <- errors.Wrap(err, "discovery error")
			continue
		}
	}
}

func (s *Service) AcquireLeadership() {
	var session *discovery.Session
	var err error
	ticker := time.NewTicker(5 * time.Second)
	for ; true; <-ticker.C {
		if session == nil {
			session, err = discovery.NewSession(s.discovery)
			if err != nil {
				s.errors <- errors.Wrap(err, "failed to register session")
				continue
			}
			session.EndlessRenew()
		}

		value, err := time.Now().UTC().MarshalJSON()
		if err != nil {
			s.errors <- errors.Wrap(err, "failed to marshal leader value")
			continue
		}

		ok, err := s.discovery.TryAcquire(&discovery.KVPair{
			Key:     s.config.LeadershipKey,
			Value:   value,
			Session: session,
		})

		if err != nil {
			s.errors <- err
			s.IsLeader = false
			continue
		}

		if ok {
			s.IsLeader = true
		} else {
			s.IsLeader = false
		}
	}
}

func (s *Service) SyncTimestamp() {
	ticker := time.NewTicker(5 * time.Second)
	s.syncTimestamp = make(chan time.Time)
	for ; true; <-ticker.C {
		kv, err := s.discovery.Get(s.config.LeadershipKey)
		if err != nil {
			s.errors <- err
			continue
		}

		if kv == nil {
			continue
		}

		t := time.Time{}
		err = t.UnmarshalJSON(kv.Value)
		if err != nil {
			s.errors <- err
			continue
		}

		s.log.Debug("timestamp synced")
		s.syncTimestamp <- t
	}
}

func (s *Service) syncTicker(interval time.Duration) chan int64 {
	ticks := make(chan int64)
	go func() {
		var cachedSyncTime time.Time
		for {
			select {
			case cachedSyncTime = <-s.syncTimestamp:
			default:
				if cachedSyncTime.IsZero() {
					time.Sleep(1 * time.Second)
					continue
				}
				now := time.Now()
				intervalsPassed := now.Sub(cachedSyncTime).Nanoseconds() / interval.Nanoseconds()
				nextSyncOffset := time.Duration((intervalsPassed + 1) * interval.Nanoseconds())
				nextSync := cachedSyncTime.Add(nextSyncOffset)
				<-time.NewTimer(nextSync.Sub(now)).C
				intervalsPassed += 1
				ticks <- intervalsPassed
			}
		}
	}()
	return ticks
}

func (s *Service) GetRates() {
	var lastTick providers.Tick
	syncTicker := s.syncTicker(5 * time.Second)

	go func() {
		// receiving tick from provider,
		// should contain latest values for all pairs
		ticks := s.provider.Ticks()
		for {
			lastTick = <-ticks
			s.log.Debug("provider ticked")
		}
	}()

	for {
		select {
		case err := <-s.provider.Errors():
			s.errors <- errors.Wrap(err, "provider error")
		case syncCount := <-syncTicker:
			tick := lastTick
			if tick == nil {
				continue
			}
			ops := tick.Ops()
			if len(ops) == 0 {
				continue
			}
			s.log.WithField("sync", syncCount).WithField("ops", len(ops)).Info("syncing")
			s.syncResults.Set(syncCount, ops)

			if !s.IsLeader {
				continue
			}

			//neighbors, err := s.discoveryService.DiscoverNeighbors()
			//if err != nil {
			//	s.errors <- err
			//	continue
			//}
			//
			//if len(neighbors) == 0 {
			//	continue
			//}

			tx := s.horizon.Transaction(&horizon.TransactionBuilder{
				Source: s.config.Master,
			})

			for _, op := range ops {
				ohaigo := op
				tx.Op(&ohaigo)
			}

			err := tx.Sign(s.config.Signer).Submit()
			if err != nil {
				s.errors <- err
				continue
			}
			s.log.WithField("sync", syncCount).Info("synced")
			//envelope, err := tx.Sign(s.config.Signer).Marshal64()
			//if err != nil {
			//	s.errors <- err
			//	continue
			//}
			//
			//b, err := json.Marshal(SyncRequest{
			//	Sync:     syncCount,
			//	Envelope: *envelope,
			//})
			//
			//if err != nil {
			//	s.errors <- err
			//	continue
			//}
			//
			//s.log.WithField("sync", syncCount).Info("attempting to verify")
			//for _, neighbor := range neighbors {
			//	fmt.Println("started")
			//	r, err := http.Post(neighbor.Address, "application/json", bytes.NewReader(b))
			//	fmt.Println("ended")
			//	if err != nil {
			//		s.errors <- err
			//		continue
			//	}
			//	if r.StatusCode == http.StatusOK {
			//		break
			//	}
			//}
			//s.log.Info("unable to verify")
		}
	}
}

func (s *Service) Run() {
	serviceSeq := []func(){
		s.Register,
		s.API,
		s.AcquireLeadership,
		s.SyncTimestamp,
		s.GetRates,
	}

	for _, fn := range serviceSeq {
		f := fn
		go func() {
			f()
		}()
	}
}

type SyncRequest struct {
	Sync     int64
	Envelope string
}
