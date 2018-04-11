package btcdeposit

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/psim/addrstate"
	"gitlab.com/swarmfund/psim/figure"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/deposits/btcdeposit/internal"
	"gitlab.com/swarmfund/psim/psim/deposits/deposit"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceBTCDeposit, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	var config Config
	err := figure.
		Out(&config).
		From(app.Config(ctx).GetRequired(conf.ServiceBTCDeposit)).
		With(figure.BaseHooks, utils.CommonHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out", logan.F{
			"service": conf.ServiceBTCDeposit,
		})
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.Signer)

	addressProvider := addrstate.New(
		ctx,
		log,
		internal.StateMutator,
		horizonConnector.Listener(),
	)

	horizonInfo, err := horizonConnector.Info()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Horizon info")
	}

	builder := xdrbuild.NewBuilder(horizonInfo.Passphrase, horizonInfo.TXExpirationPeriod)

	return deposit.New(&deposit.Opts{
		log,
		config.Source,
		config.Signer,
		conf.ServiceBTCDeposit,
		conf.ServiceBTCDepositVerify,

		config.LastProcessedBlock,
		config.LastBlocksNotWatch,

		horizonConnector,
		addressProvider,
		globalConfig.Discovery(),
		builder,
		NewBTCHelper(
			log,

			config.DepositAsset,
			config.MinDepositAmount,
			config.FixedDepositFee,

			globalConfig.Bitcoin(),
		),
	}), nil
}
