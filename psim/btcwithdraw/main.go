package btcwithdraw

import (
	"context"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/swarmfund/psim/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
)

func init() {
	app.RegisterService(conf.ServiceBTCWithdraw, setupFn)
}

func setupFn(ctx context.Context) (utils.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	var config Config
	err := figure.
		Out(&config).
		From(app.Config(ctx).Get(conf.ServiceBTCWithdraw)).
		With(figure.BaseHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out", logan.F{
			"service": conf.ServiceBTCWithdraw,
		})
	}

	horizonConnector, err := globalConfig.Horizon()
	if err != nil {
		panic(err)
	}

	return New(log, config, globalConfig.HorizonV2().Listener(), horizonConnector, nil), nil
}
