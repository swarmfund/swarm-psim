package dashwithdraw

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/tokend/go/xdrbuild"
)

func init() {
	app.RegisterService(conf.ServiceDashWithdraw, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	config, err := NewConfig(app.Config(ctx).GetRequired(conf.ServiceDashWithdraw))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create config", logan.F{
			"service": conf.ServiceDashWithdraw,
		})
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.SignerKP)

	horizonInfo, err := horizonConnector.System().Info()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Horizon info")
	}

	builder := xdrbuild.NewBuilder(horizonInfo.Passphrase, horizonInfo.TXExpirationPeriod)
	btcHelper, err := NewDashHelper(
		log,
		config,
		globalConfig.Bitcoin(),
		NewRandomCoinSelector(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create CommonDashHelper")
	}

	// FIXME
	// FIXME
	// FIXME
	//return withdraw.New(
	//	conf.ServiceDashWithdraw,
	//	config.SignerKP,
	//	log,
	//	horizonConnector.Listener(),
	//	horizonConnector.Operations(),
	//	horizonConnector.Submitter(),
	//	builder,
	//	withdraw.VerificationConfig{
	//		Verify:   false,
	//		SourceKP: config.SourceKP,
	//	},
	//	btcHelper,
	//), nil

	builder = builder
	return btcHelper, nil
}
