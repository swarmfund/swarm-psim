package btcdepositveri

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/addrstate"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/deposits/btcdeposit"
	"gitlab.com/swarmfund/psim/psim/deposits/depositveri"
)

func init() {
	app.RegisterService(conf.ServiceBTCDepositVerify, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	config, err := NewConfig(globalConfig.GetRequired(conf.ServiceBTCDepositVerify))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create config", logan.F{
			"service": conf.ServiceBTCDepositVerify,
		})
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.Signer)

	btcHelper, err := btcdeposit.NewBTCHelper(
		log,
		config.DepositAsset,
		config.MinDepositAmount,
		config.FixedDepositFee,
		config.OffchainCurrency,
		config.OffchainBlockchain,
		config.BlocksToSearchForTX,
		globalConfig.Bitcoin(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create BTCHelper")
	}

	addrProvider := addrstate.New(
		ctx,
		app.Log(ctx),
		[]addrstate.StateMutator{
			addrstate.ExternalSystemBindingMutator{SystemType: config.ExternalSystem},
			addrstate.BalanceMutator{Asset: config.DepositAsset},
		},
		horizonConnector.Listener(),
	)

	builder, err := horizonConnector.TXBuilder()
	if err != nil {
		return nil, errors.Wrap(err, "failed to init tx builder")
	}

	return depositveri.New(depositveri.Opts{
		Log:                log.WithField("service", conf.ServiceBTCDepositVerify),
		Source:             config.Source,
		Signer:             config.Signer,
		ExternalSystem:     config.ExternalSystem,
		LastBlocksNotWatch: config.LastBlocksNotWatch,
		Horizon:            horizonConnector,
		IssuanceStreamer:   horizonConnector.Listener(),
		AddressProvider:    addrProvider,
		Builder:            builder,
		OffchainHelper:     btcHelper,
	}), nil
}
