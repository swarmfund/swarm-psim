package wallet_cleaner

import (
	"context"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
)

func init() {
	app.RegisterService(conf.ServiceWalletCleaner, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)

	config := Config{}
	err := figure.Out(&config).From(globalConfig.Get(conf.ServiceWalletCleaner)).Please()
	if err != nil {
		return nil, errors.Wrap(err, "failed to figure out")
	}

	return New(app.Log(ctx), app.Config(ctx).Horizon(), config.ExpireDuration), nil
}
