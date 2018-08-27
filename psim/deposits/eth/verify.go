package eth

import (
	"context"

	"gitlab.com/swarmfund/psim/psim/deposits/depositveri"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	internal2 "gitlab.com/swarmfund/psim/psim/deposits/eth/internal"
	"gitlab.com/swarmfund/psim/psim/internal"
	"gitlab.com/tokend/addrstate"
)

func init() {
	app.RegisterService(conf.ServiceETHDepositVerify, func(ctx context.Context) (app.Service, error) {
		config, err := NewDepositVerifyConfig(app.Config(ctx).Get(conf.ServiceETHDepositVerify))
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

		return depositveri.New(depositveri.Opts{
			Log:                app.Log(ctx),
			Source:             config.Source,
			Signer:             config.Signer,
			ExternalSystem:     config.ExternalSystem,
			LastBlocksNotWatch: config.Confirmations,
			Horizon:            hrz,
			IssuanceStreamer:   hrz.Listener(),
			AddressProvider:    addressProvider,
			Builder:            txbuilder,
			OffchainHelper:     helper,
		}), nil
	})
}
