// BTC Withdraw Verify
package btcwithdveri

import (
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/swarmfund/psim/figure"
	"context"
	"gitlab.com/swarmfund/psim/ape"
)

func init() {
	app.RegisterService(conf.ServiceBTCWithdrawVerify, setupFn)
}

func setupFn(ctx context.Context) (utils.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	var config Config
	err := figure.
		Out(&config).
		From(app.Config(ctx).GetRequired(conf.ServiceBTCWithdrawVerify)).
		With(figure.BaseHooks, utils.CommonHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out", logan.F{
			"service": conf.ServiceBTCWithdrawVerify,
		})
	}

	horizonConnector, err := globalConfig.Horizon()
	if err != nil {
		panic(err)
	}

	listener, err := ape.Listener(config.Host, config.Port)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init listener")
	}

	return New(log, config, horizonConnector, globalConfig.Bitcoin(), listener), nil
}
