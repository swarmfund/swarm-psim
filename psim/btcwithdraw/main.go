package btcwithdraw

import (
	"context"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceBTCWithdraw, setupFn)
}

func setupFn(ctx context.Context) (utils.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	horizonConnector, err := globalConfig.Horizon()
	if err != nil {
		panic(err)
	}

	return New(log, globalConfig.HorizonV2().Listener(), horizonConnector, nil), nil
}
