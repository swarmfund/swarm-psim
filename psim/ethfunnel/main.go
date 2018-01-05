package ethfunnel

import (
	"context"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/ethfunnel/internal"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceETHFunnel, func(ctx context.Context) (utils.Service, error) {
		config := Config{}
		err := figure.
			Out(&config).
			From(app.Config(ctx).GetRequired(conf.ServiceETHFunnel)).
			With(figure.BaseHooks, utils.ETHHooks).
			Please()
		if err != nil {
			return nil, errors.Wrap(err, "failed to figure out")
		}

		eth := app.Config(ctx).Ethereum()

		wallet, err := internal.NewWallet(config.Seed)
		if err != nil {
			return nil, errors.Wrap(err, "failed to init wallet")
		}

		return NewService(ctx, app.Log(ctx), config, wallet, eth), nil
	})
}
