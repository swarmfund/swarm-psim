package erc20

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/deposits/depositveri"
	"gitlab.com/swarmfund/psim/psim/deposits/erc20/internal"
	. "gitlab.com/swarmfund/psim/psim/internal"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/addrstate"
)

func init() {
	app.RegisterService(conf.ServiceERC20DepositVerify, func(ctx context.Context) (app.Service, error) {
		config := VerifyConfig{
			Confirmations: 12,
		}
		err := figure.
			Out(&config).
			With(figure.BaseHooks, utils.ETHHooks).
			From(app.Config(ctx).Get(conf.ServiceERC20DepositVerify)).
			Please()
		if err != nil {
			return nil, errors.Wrap(err, "failed to figure out")
		}

		horizon := app.Config(ctx).Horizon().WithSigner(config.Signer)

		if config.ExternalSystem == 0 {
			config.ExternalSystem = MustGetExternalSystemType(horizon.Assets(), config.DepositAsset)
		}

		addrProvider := addrstate.New(
			ctx,
			app.Log(ctx),
			[]addrstate.StateMutator{
				addrstate.ExternalSystemBindingMutator{SystemType: config.ExternalSystem},
				addrstate.BalanceMutator{Asset: config.DepositAsset},
			},
			horizon.Listener(),
		)

		builder, err := horizon.TXBuilder()
		if err != nil {
			return nil, errors.Wrap(err, "failed to init tx builder")
		}

		eth := app.Config(ctx).Ethereum()

		return depositveri.New(depositveri.Opts{
			Log:                app.Log(ctx).WithField("service", conf.ServiceERC20DepositVerify),
			Source:             config.Source,
			Signer:             config.Signer,
			ExternalSystem:     config.ExternalSystem,
			LastBlocksNotWatch: config.Confirmations,
			Horizon:            horizon,
			IssuanceStreamer:   horizon.Listener(),
			AddressProvider:    addrProvider,
			Builder:            builder,
			OffchainHelper:     internal.NewERC20Helper(eth, config.DepositAsset, config.Token, config.BlocksToSearchForTX),
		}), nil
	})
}
