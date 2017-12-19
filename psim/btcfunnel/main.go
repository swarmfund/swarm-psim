package btcfunnel

import (
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/utils"
	"context"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/figure"
	"fmt"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func init() {
	app.RegisterService(conf.ServiceBTCFunnel, setupFn)
}

func setupFn(ctx context.Context) (utils.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx).WithField("service", conf.ServiceBTCFunnel)

	config := Config{}

	err := figure.
		Out(&config).
		From(globalConfig.Get(conf.ServiceBTCFunnel)).
		With(figure.BaseHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to figure out %s", conf.ServiceBTCFunnel))
	}

	btcClient, err := globalConfig.Bitcoin()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Bitcoin client from global config")
	}

	return New(config, log, btcClient), nil
}
