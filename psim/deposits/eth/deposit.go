package eth

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/deposits/deposit"
	internal2 "gitlab.com/swarmfund/psim/psim/deposits/eth/internal"
	"gitlab.com/swarmfund/psim/psim/internal"
	"gitlab.com/tokend/addrstate"
)

func init() {
	app.RegisterService(conf.ServiceETHDeposit, func(ctx context.Context) (app.Service, error) {
		config, err := NewDepositConfig(app.Config(ctx).Get(conf.ServiceETHDeposit))
		if err != nil {
			return nil, errors.Wrap(err, "failed to init config")
		}

		hrz := app.Config(ctx).Horizon().WithSigner(config.Signer)

		if config.ExternalSystem == 0 {
			config.ExternalSystem = internal.MustGetExternalSystemType(hrz.Assets(), config.DepositAsset)
		}

		addressProvider := addrstate.New(
			ctx,
			app.Log(ctx),
			[]addrstate.StateMutator{
				addrstate.ExternalSystemBindingMutator{SystemType: config.ExternalSystem},
				addrstate.BalanceMutator{Asset: config.DepositAsset},
			},
			hrz.Listener(),
		)

		txbuilder, err := hrz.TXBuilder()
		if err != nil {
			return nil, errors.Wrap(err, "failed to init tx builder")
		}

		helper := internal2.ETHHelper{
			config.DepositAsset,
			config.MinDepositAmount,
			config.FixedDepositFee,
			config.BlocksToSearchForTX,
			app.Config(ctx).Ethereum(),
		}

		return deposit.New(&deposit.Opts{
			Log:                app.Log(ctx),
			Source:             config.Source,
			Signer:             config.Signer,
			ServiceName:        conf.ServiceETHDeposit,
			LastProcessedBlock: config.Cursor,
			Horizon:            hrz,
			ExternalSystem:     config.ExternalSystem,
			AddressProvider:    addressProvider,
			Builder:            txbuilder,
			OffchainHelper:     helper,
		}), nil
	})
}
