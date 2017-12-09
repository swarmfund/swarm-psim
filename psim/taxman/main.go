package taxman

import (
	"context"
	"fmt"

	api "gitlab.com/tokend/psim/ape"
	"gitlab.com/tokend/psim/figure"
	"gitlab.com/tokend/psim/psim/app"
	"gitlab.com/tokend/psim/psim/conf"
	"gitlab.com/tokend/psim/psim/utils"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func init() {
	app.RegisterService(conf.ServiceTaxman, func(ctx context.Context) (utils.Service, error) {
		serviceConfig := Config{
			ServiceName:   conf.ServiceTaxman,
			Host:          "localhost",
			LeadershipKey: fmt.Sprintf("service/%s/leader", conf.ServiceTaxman),
		}

		globalConfig := ctx.Value(app.CtxConfig).(conf.Config)
		err := figure.Out(&serviceConfig).
			From(globalConfig.Get(conf.ServiceTaxman)).
			With(figure.BaseHooks, utils.CommonHooks, SkipTransactionsHook).
			Please()
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to figure out %s", conf.ServiceTaxman))
		}

		log := ctx.Value(app.CtxLog).(*logan.Entry)

		discovery, err := globalConfig.Discovery()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get discovery client")
		}

		horizon, err := globalConfig.Horizon()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get horizon connector")
		}

		listener, err := api.Listener(serviceConfig.Host, serviceConfig.Port)
		if err != nil {
			return nil, errors.Wrap(err, "failed to init listener")
		}

		log.Info("starting")

		service := New(log, discovery, horizon, listener, serviceConfig, ctx)
		return service, nil
	})
}
