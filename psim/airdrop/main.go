package airdrop

import (
	"context"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
)

func init() {
	app.RegisterService(conf.ServiceAirdrop, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	var config Config
	err := figure.
		Out(&config).
		From(app.Config(ctx).GetRequired(conf.ServiceRateSync)).
		With(figure.BaseHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out", logan.F{
			"service": conf.ServiceRateSync,
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
		builder,
		horizonConnector.Submitter(),
		config.Source,
		config.Signer,
	), nil
}
