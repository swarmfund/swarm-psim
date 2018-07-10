package eventsubmitter

import (
	"context"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.EventSubmitterService, setupService)
}

func setupService(ctx context.Context) (app.Service, error) {
	var serviceConfig ServiceConfig
	serviceConfigMap := app.Config(ctx).GetRequired(conf.EventSubmitterService)

	err := figure.Out(&serviceConfig).From(serviceConfigMap).With(figure.BaseHooks, utils.ETHHooks).Please()
	if err != nil {
		return nil, errors.Wrap(err, "failed to figure-out")
	}

	logger := app.Log(ctx)

	horizonConnector := app.Config(ctx).Horizon().WithSigner(serviceConfig.Signer)
	extractor := NewTokendExtractor(logger, horizonConnector.Listener().StreamTXsFromCursor(ctx, serviceConfig.TxHistoryCursor, false))
	handler := NewTokendHandler(logger, horizonConnector).withTokendProcessors()
	broadcaster := NewGenericBroadcaster(logger)

	broadcaster.AddTarget(NewSalesforceTarget(app.Config(ctx).Salesforce()))
	broadcaster.AddTarget(NewMixpanelTarget(app.Config(ctx).Mixpanel()))

	return NewService(serviceConfig, extractor, handler, broadcaster, logger), nil
}
