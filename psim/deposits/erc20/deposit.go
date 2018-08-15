package erc20

import (
	"context"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/addrstate"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/deposits/deposit"
	"gitlab.com/swarmfund/psim/psim/deposits/erc20/internal"
	. "gitlab.com/swarmfund/psim/psim/internal"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceERC20Deposit, func(ctx context.Context) (app.Service, error) {
		var config DepositConfig

		err := figure.
			Out(&config).
			With(figure.BaseHooks, utils.ETHHooks).
			From(app.Config(ctx).Get(conf.ServiceERC20Deposit)).
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

		return deposit.New(&deposit.Opts{
			Log:                 app.Log(ctx),
			Source:              config.Source,
			Signer:              config.Signer,
			ServiceName:         conf.ServiceERC20Deposit,
			VerifierServiceName: conf.ServiceERC20DepositVerify,
			LastProcessedBlock:  config.Cursor,
			Horizon:             app.Config(ctx).Horizon().WithSigner(config.Signer),
			ExternalSystem:      config.ExternalSystem,
			AddressProvider:     addrProvider,
			Discovery:           app.Config(ctx).Discovery(),
			Builder:             builder,
			OffchainHelper:      internal.NewERC20Helper(eth, config.DepositAsset, config.Token),
			DisableVerify:       config.DisableVerify,
		}), nil
	})
}
