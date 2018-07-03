package reporter

import (
	"context"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.BalancesReporterService, setupService)
}

func setupService(ctx context.Context) (app.Service, error) {
	var serviceConfig ServiceConfig
	serviceConfigMap := app.Config(ctx).GetRequired(conf.BalancesReporterService)

	err := figure.Out(&serviceConfig).From(serviceConfigMap).With(figure.BaseHooks, utils.ETHHooks).Please()
	if err != nil {
		return nil, errors.Wrap(err, "failed to figure-out")
	}

	logger := app.Log(ctx)
	horizon := app.Config(ctx).Horizon().WithSigner(serviceConfig.Signer)
	broadcaster := NewGenericBroadcaster(logger)
	broadcaster.AddTarget(NewSalesforceTarget(app.Config(ctx).Salesforce()))

	return NewService(serviceConfig, horizon, broadcaster, logger), nil
}
