package btcdepositveri

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/figure"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/btcdeposit"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/depositveri"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceBTCDepositVerify, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	var config Config
	err := figure.
		Out(&config).
		From(app.Config(ctx).GetRequired(conf.ServiceBTCDepositVerify)).
		With(figure.BaseHooks, utils.CommonHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out", logan.F{
			"service": conf.ServiceBTCDepositVerify,
		})
	}

	listener, err := ape.Listener(config.Host, config.Port)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init listener")
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.Signer)

	horizonInfo, err := horizonConnector.Info()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Horizon info")
	}

	return depositveri.New(
		"bitcoin",
		conf.ServiceBTCDepositVerify,
		log,
		config.Signer,
		config.LastBlocksNotWatch,
		horizonConnector,
		xdrbuild.NewBuilder(horizonInfo.Passphrase, horizonInfo.TXExpirationPeriod),
		listener,
		globalConfig.Discovery(),
		btcdeposit.NewBTCHelper(
			log,
			config.DepositAsset,
			config.MinDepositAmount,
			config.FixedDepositFee,
			globalConfig.Bitcoin(),
		),
	), nil
}
