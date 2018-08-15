package btcdeposit

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/addrstate"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/deposits/deposit"
	"gitlab.com/tokend/go/xdrbuild"
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
		[]addrstate.StateMutator{
			addrstate.ExternalSystemBindingMutator{SystemType: config.ExternalSystem},
			addrstate.BalanceMutator{Asset: config.DepositAsset},
		},
		horizonConnector.Listener(),
	)

	horizonInfo, err := horizonConnector.System().Info()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Horizon info")
	}

	builder := xdrbuild.NewBuilder(horizonInfo.Passphrase, horizonInfo.TXExpirationPeriod)
	btcHelper, err := NewBTCHelper(
		log,

		config.DepositAsset,
		config.MinDepositAmount,
		config.FixedDepositFee,
		config.OffchainCurrency,
		config.OffchainBlockchain,

		globalConfig.Bitcoin(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create BTCHelper")
	}

	return deposit.New(&deposit.Opts{
		Log:                 log,
		Source:              config.Source,
		Signer:              config.Signer,
		ServiceName:         conf.ServiceBTCDeposit,
		VerifierServiceName: conf.ServiceBTCDepositVerify,
		LastProcessedBlock:  config.LastProcessedBlock,
		Horizon:             horizonConnector,
		ExternalSystem:      config.ExternalSystem,
		AddressProvider:     addressProvider,
		Discovery:           globalConfig.Discovery(),
		Builder:             builder,
		OffchainHelper:      btcHelper,
		DisableVerify:       config.DisableVerify,
	}), nil
}
