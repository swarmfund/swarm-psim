package charger

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/stripe/stripe-go/client"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/figure"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
)

func init() {
	app.Register(conf.ServiceCharger, func(ctx context.Context) error {
		serviceConfig := Config{
			Host: "localhost",
		}

		globalConfig := ctx.Value(app.CtxConfig).(conf.Config)
		err := figure.
			Out(&serviceConfig).
			From(globalConfig.Get(conf.ServiceCharger)).
			With(figure.BaseHooks).
			Please()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to figure out %s", conf.ServiceCharger))
		}

		log := ctx.Value(app.CtxLog).(*logan.Entry)

		stripe, err := globalConfig.Stripe()
		if err != nil {
			return errors.Wrap(err, "failed to init stripe client")
		}

		listener, err := ape.Listener(serviceConfig.Host, serviceConfig.Port)
		if err != nil {
			return errors.Wrap(err, "failed to init listener")
		}

		retryTicker := time.NewTicker(10 * time.Second)

		service, errs := New(log, serviceConfig, stripe, listener)
		for service.Run(); ; <-retryTicker.C {
			log.WithError(<-errs).Warn("service error")
		}
	})
}

type Service struct {
	log      *logan.Entry
	config   Config
	stripe   *client.API
	listener net.Listener
	errors   chan error
}

func New(
	log *logan.Entry, config Config, stripe *client.API, listener net.Listener,
) (*Service, chan error) {
	service := &Service{
		log:      log,
		config:   config,
		stripe:   stripe,
		listener: listener,
		errors:   make(chan error),
	}
	return service, service.errors
}
