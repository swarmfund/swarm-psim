package ordernotifier

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/figure"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceOrderNotifier, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	var config Config
	err := figure.
		Out(&config).
		From(globalConfig.GetRequired(conf.ServiceOrderNotifier)).
		With(figure.BaseHooks, utils.CommonHooks, EmailsHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "failed to figure out", logan.F{
			"service": conf.ServiceOrderNotifier,
		})
	}

	if len(config.RequestTokenSuffix) == 0 {
		return nil, errors.New("'email_request_token_suffix' in config must not be empty")
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.Signer)

	checkSaleStateResponses := horizonConnector.Listener().StreamAllCheckSaleStateOps(ctx, 0)

	return New(
		config,
		horizonConnector.Transactions(),
		globalConfig.Notificator(),
		log,
		horizonConnector.Users(),
		horizonConnector.Sales(),
		checkSaleStateResponses,
	), nil
}
