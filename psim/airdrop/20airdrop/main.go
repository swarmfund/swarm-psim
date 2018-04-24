package mrefairdrop

import (
	"context"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceAirdropMarch20b20, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	var config Config
	err := figure.
		Out(&config).
		From(app.Config(ctx).GetRequired(conf.ServiceAirdropMarch20b20)).
		With(figure.BaseHooks, utils.ETHHooks, airdrop.FigureHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out", logan.F{
			"service": conf.ServiceAirdropMarch20b20,
		})
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.Signer)

	horizonInfo, err := horizonConnector.Info()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Horizon info")
	}

	builder := xdrbuild.NewBuilder(horizonInfo.Passphrase, horizonInfo.TXExpirationPeriod)

	issuanceSubmitter := airdrop.NewIssuanceSubmitter(
		config.IssuanceAsset,
		airdrop.March20b20ReferenceSuffix,
		config.Source,
		config.Signer,
		builder,
		horizonConnector.Submitter())

	ledgerStreamer := airdrop.NewLedgerChangesStreamer(log, horizonConnector.Listener())

	emailProcessor := airdrop.NewEmailsProcessor(log, config.EmailsConfig, globalConfig.Notificator())

	return NewService(
		log,
		config,
		issuanceSubmitter,
		ledgerStreamer,
		airdrop.NewBalanceIDProvider(horizonConnector.Accounts()),
		horizonConnector.Users(),
		emailProcessor,
	), nil
}
