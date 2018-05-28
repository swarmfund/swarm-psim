package listener

import (
	"context"

	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/swarmfund/psim/psim/listener/internal"
)

func init() {
	app.RegisterService(conf.ListenerService, setupService)
}

// TODO make targets to be loaded from config
func setupService(ctx context.Context) (app.Service, error) {
	appConfig := app.Config(ctx)
	listenerConfig := ServiceConfig{}
	err := figure.
		Out(&listenerConfig).
		From(appConfig.Get(conf.ListenerService)).
		With(figure.BaseHooks, utils.ETHHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "failed to figure out extractor config")
	}
	horizonConnector := appConfig.Horizon().WithSigner(listenerConfig.Signer)
	var txStreamResponse TokendExtractor = horizonConnector.Listener().StreamTXsFromCursor(ctx, "", false)
	logger := app.Log(ctx).WithField("service", conf.ListenerService)
	return NewService(listenerConfig, txStreamResponse, *NewTokendHandler(horizonConnector.Operations(), horizonConnector.Accounts()), internal.NewGenericBroadcaster(NewMixpanelTarget("371c0e972d45715b770c176386f82971")), logger), nil
}
