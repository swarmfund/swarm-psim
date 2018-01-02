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
		From(globalConfig.GetRequired(conf.ServiceBTCFunnel)).
		With(figure.BaseHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to figure out %s", conf.ServiceBTCFunnel))
	}

	// TODO Validate config. Some values can't be zero.

	return New(config, log, globalConfig.Bitcoin()), nil
}
