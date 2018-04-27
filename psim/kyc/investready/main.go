package investready

import (
	"context"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/swarmfund/psim/psim/kyc"
)

func init() {
	app.RegisterService(conf.ServiceInvestReady, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	var config Config
	err := figure.
		Out(&config).
		From(app.Config(ctx).GetRequired(conf.ServiceInvestReady)).
		With(figure.BaseHooks, utils.ETHHooks, hooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out", logan.F{
			"service": conf.ServiceInvestReady,
		})
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.Signer)

	horizonInfo, err := horizonConnector.Info()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Horizon info")
	}

	builder := xdrbuild.NewBuilder(horizonInfo.Passphrase, horizonInfo.TXExpirationPeriod)

	return NewService(
		log,
		config,
		horizonConnector.Listener(),
		horizonConnector.Operations(),
		kyc.NewRequestPerformer(builder, config.Source, config.Signer, horizonConnector.Submitter()),
		horizonConnector.Blobs(),
		kyc.NewBlobDataRetriever(horizonConnector.Blobs()),
		horizonConnector.Users(),
		NewConnector(log.WithField("service", conf.ServiceInvestReady), config.Connector),
	), nil
}
