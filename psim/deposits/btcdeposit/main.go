package btcdeposit

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/deposits/deposit"
	"gitlab.com/swarmfund/psim/psim/internal"
	"gitlab.com/tokend/addrstate"
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

	if config.ExternalSystem == 0 {
		config.ExternalSystem = internal.MustGetExternalSystemType(horizonConnector.Assets(), config.DepositAsset)
	}

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
		config.NetworkType,
		// this value is actually only needed for btcdepositveri service
		10,

		globalConfig.Bitcoin(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create BTCHelper")
	}

	return deposit.New(&deposit.Opts{
		Log:                log,
		Source:             config.Source,
		Signer:             config.Signer,
		ServiceName:        conf.ServiceBTCDeposit,
		LastProcessedBlock: config.LastProcessedBlock,
		Horizon:            horizonConnector,
		ExternalSystem:     config.ExternalSystem,
		AddressProvider:    addressProvider,
		Builder:            builder,
		OffchainHelper:     btcHelper,
	}), nil
}
