package listener

import (
	"context"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/listener/internal"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
)

// TODO move salesforce and mixpanel

func init() {
	app.RegisterService(conf.ListenerService, setupService)
}

func setupService(ctx context.Context) (app.Service, error) {
	serviceConfig, err := loadServiceConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load service config")
	}

	logger := app.Log(ctx).WithField("service", conf.ListenerService)

	//app.Config(ctx).Salesforce().WithUserData(username, password)

	horizonConnector := getHorizonConnector(ctx, serviceConfig)
	extractor := getExtractor(ctx, horizonConnector, serviceConfig.TxhistoryCursor)
	handler := setupHandler(horizonConnector)
	targets := setupBroadcasterTargets(serviceConfig.MixpanelToken, serviceConfig.SalesforceUsername, serviceConfig.SalesforcePassword)
	broadcaster, err := setupBroadcaster(targets)

	if err != nil {
		return nil, errors.Wrap(err, "failed to setup broadcaster")
	}

	return NewService(serviceConfig, extractor, handler, broadcaster, logger), nil
}

func loadServiceConfig(ctx context.Context) (ServiceConfig, error) {
	serviceConfig := ServiceConfig{}
	serviceConfigMap := app.Config(ctx).GetRequired(conf.ListenerService)

	err := figure.Out(&serviceConfig).From(serviceConfigMap).With(figure.BaseHooks, utils.ETHHooks).Please()
	if err != nil {
		return ServiceConfig{}, errors.Wrap(err, "failed to parse service config using 'figure-out' from map")
	}

	return serviceConfig, nil
}

func getExtractor(ctx context.Context, connector horizon.Connector, txhistoryCursor string) TokendExtractor {
	return connector.Listener().StreamTXsFromCursor(ctx, txhistoryCursor, false)
}

func setupHandler(connector horizon.Connector) TokendHandler {
	return *NewTokendHandler().withTokendProcessors(connector)
}

func getHorizonConnector(ctx context.Context, config ServiceConfig) horizon.Connector {
	return *app.Config(ctx).Horizon().WithSigner(config.Signer)
}

func (th *TokendHandler) withTokendProcessors(connector horizon.Connector) *TokendHandler {
	th.HorizonConnector = connector
	th.SetProcessor(xdr.OperationTypeCreateKycRequest, processKYCCreateUpdateRequestOp)
	th.SetProcessor(xdr.OperationTypeReviewRequest, processReviewRequestOp(connector.Operations(), connector.Accounts()))
	th.SetProcessor(xdr.OperationTypeCreateIssuanceRequest, processCreateIssuanceRequestOp)
	th.SetProcessor(xdr.OperationTypeManageOffer, processManageOfferOp)
	th.SetProcessor(xdr.OperationTypePayment, processPayment)
	th.SetProcessor(xdr.OperationTypePaymentV2, processPaymentV2)
	th.SetProcessor(xdr.OperationTypeCreateWithdrawalRequest, processWithdrawRequest)
	th.SetProcessor(xdr.OperationTypeCreateAccount, processCreateAccountOp)
	return th
}

// MaybeTarget can contain actual Target or error
type MaybeTarget struct {
	Target
	error
}

func setupBroadcasterTargets(mixpanelToken string, salesforceUsername string, salesforcePassword string) (targets []MaybeTarget) {
	salesforceConnector := NewSalesforceConnector(app.Config(nil).Salesforce().WithUserData(salesforceUsername, salesforcePassword))
	salesforceTarget := salesforceConnector.GetTarget()
	targets = append(targets, MaybeTarget{salesforceTarget, nil})
	mixpanelTarget := NewMixpanelTarget(mixpanelToken)
	targets = append(targets, MaybeTarget{mixpanelTarget, nil})
	return targets
}

func setupBroadcaster(maybeTargets []MaybeTarget) (*GenericBroadcaster, error) {
	broadcaster := internal.NewGenericBroadcaster()
	for _, target := range maybeTargets {
		if target.error != nil {
			return nil, errors.Wrap(target.error, "invalid target received")
		}
		broadcaster.AddTarget(target)
	}
	return broadcaster, nil
}
