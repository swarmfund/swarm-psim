package telegram

import (
	"context"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/listener"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/go/xdrbuild"
)

func init() {
	app.RegisterService(conf.ServiceAirdropTelegram, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	var config Config
	err := figure.
		Out(&config).
		From(app.Config(ctx).GetRequired(conf.ServiceAirdropTelegram)).
		With(figure.BaseHooks, utils.ETHHooks, airdrop.FigureHooks, listener.ConfigFigureHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out", logan.F{
			"service": conf.ServiceAirdropTelegram,
		})
	}

	validateErr := config.Validate()
	if validateErr != nil {
		return nil, errors.Wrap(validateErr, "Config is invalid")
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.Signer)

	horizonInfo, err := horizonConnector.System().Info()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Horizon info")
	}

	builder := xdrbuild.NewBuilder(horizonInfo.Passphrase, horizonInfo.TXExpirationPeriod)

	issuanceSubmitter := airdrop.NewIssuanceSubmitter(
		config.Issuance.Asset,
		airdrop.TelegramReferenceSuffix,
		config.Source,
		config.Signer,
		builder,
		horizonConnector.Submitter())

	return NewService(
		log,
		config,
		issuanceSubmitter,
		airdrop.NewBalanceIDProvider(horizonConnector.Accounts()),
	), nil
}
