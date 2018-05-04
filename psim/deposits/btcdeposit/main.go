package btcdeposit

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/swarmfund/psim/addrstate"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/deposits/btcdeposit/internal"
	"gitlab.com/swarmfund/psim/psim/deposits/deposit"
)

func init() {
	app.RegisterService(conf.ServiceBTCDeposit, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	config, err := NewConfig(globalConfig.GetRequired(conf.ServiceBTCDeposit))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create config", logan.F{
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
