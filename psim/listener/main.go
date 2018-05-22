package listener

import (
	"context"

	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func init() {
	app.RegisterService(conf.ListenerService, setup)
}

func setup(ctx context.Context) (app.Service, error) {
	appConfig := app.Config(ctx)
	listenerConfig := Config{}

	err := figure.
		Out(&listenerConfig).
		From(appConfig.Get(conf.ListenerService)).
		With(figure.BaseHooks, utils.ETHHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "failed to figure out listener config")
	}
	horizonConnector := appConfig.Horizon().WithSigner(listenerConfig.Signer)
	txStreamResponse := horizonConnector.Listener().StreamTXsFromCursor(ctx, "now", false)
	requestsProvider := horizonConnector.Operations()
	accountsProvider := horizonConnector.Accounts()
	logger := app.Log(ctx).WithField("service", conf.ListenerService)
	listener := NewListener(requestsProvider, txStreamResponse, accountsProvider, logger)
	return New(listenerConfig, *listener, logger), nil
}
