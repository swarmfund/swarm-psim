package kycairdrop

import (
	"context"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
)

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

	if len(config.RequestTokenSuffix) == 0 {
		return nil, errors.New("'email_request_token_suffix' in config must not be empty")
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.Signer)

	horizonInfo, err := horizonConnector.Info()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Horizon info")
	}

	builder := xdrbuild.NewBuilder(horizonInfo.Passphrase, horizonInfo.TXExpirationPeriod)

	issuanceSubmitter := airdrop.NewIssuanceSubmitter(
		config.Asset,
		airdrop.KYCReferenceSuffix,
		config.Source,
		config.Signer,
		builder,
		horizonConnector.Submitter())

	emailProcessor := airdrop.NewEmailsProcessor(log, config.EmailsConfig, globalConfig.Notificator())

	ledgerChangesStreamer := airdrop.NewLedgerChangesStreamer(log, horizonConnector.Listener())

	return NewService(
		log,
		config,
		issuanceSubmitter,
		ledgerChangesStreamer,
		horizonConnector.Accounts(),
		horizonConnector.Users(),
		horizonConnector.Accounts(),
		emailProcessor,
	), nil
}
