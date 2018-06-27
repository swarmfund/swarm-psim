package telegram

import (
	"context"

	"github.com/andstepko/mtproto"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/listener"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/go/doorman"
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

	err = config.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "Config is invalid")
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.Signer)

	builder, err := horizonConnector.TXBuilder()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Horizon TXBuilder")
	}

	issuanceSubmitter := airdrop.NewIssuanceSubmitter(
		config.Issuance.Asset,
		airdrop.TelegramReferenceSuffix,
		config.Source,
		config.Signer,
		builder,
		horizonConnector.Submitter())

	storage, err := mtproto.NewStringSecretsStorage(config.TelegramSecretKey)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create Telegram secrets storage from hex string")
	}

	telegram, err := mtproto.NewMTProto(storage)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create Telegram connector")
	}
	err = telegram.Connect()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to Telegram")
	}

	return NewService(
		log,
		config,
		telegram,
		issuanceSubmitter,
		airdrop.NewBalanceIDProvider(horizonConnector.Accounts()),
		doorman.New(!config.Listener.CheckSignature, horizonConnector.Accounts()),
	), nil
}
