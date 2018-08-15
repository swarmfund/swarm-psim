package kycairdrop

import (
	"context"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/lchanges"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/go/xdrbuild"
)

// KYCAirdrop service's goal is to issue 10 SWM to users who have successfully passed KYC.

func init() {
	app.RegisterService(conf.ServiceAirdropKYC, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	var config Config
	err := figure.
		Out(&config).
		From(app.Config(ctx).GetRequired(conf.ServiceAirdropKYC)).
		With(figure.BaseHooks, utils.ETHHooks, airdrop.FigureHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out", logan.F{
			"service": conf.ServiceAirdropKYC,
		})
	}

	err = config.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "Config is invalid")
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.Signer)

	horizonInfo, err := horizonConnector.System().Info()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Horizon info")
	}

	builder := xdrbuild.NewBuilder(horizonInfo.Passphrase, horizonInfo.TXExpirationPeriod)

	if config.IssuanceConfig.ReferenceSuffix == "" {
		config.IssuanceConfig.ReferenceSuffix = airdrop.KYCReferenceSuffix
	}

	issuanceSubmitter := airdrop.NewIssuanceSubmitter(
		config.IssuanceConfig.Asset,
		config.IssuanceConfig.ReferenceSuffix,
		config.Source,
		config.Signer,
		builder,
		horizonConnector.Submitter())

	emailProcessor := airdrop.NewEmailsProcessor(log, config.EmailsConfig, globalConfig.Notificator())

	ledgerChangesStreamer := lchanges.NewStreamer(log, horizonConnector.Listener(), false)

	return NewService(
		log,
		config,
		issuanceSubmitter,
		ledgerChangesStreamer,
		horizonConnector.Accounts(),
		horizonConnector.Users(),
		airdrop.NewUSAChecker(horizonConnector.Blobs()),
		horizonConnector.Accounts(),
		emailProcessor,
	), nil
}
