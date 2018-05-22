package erc20

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/psim/addrstate"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/deposits/deposit"
	"gitlab.com/swarmfund/psim/psim/deposits/erc20/internal"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceERC20Deposit, func(ctx context.Context) (app.Service, error) {
		config := DepositConfig{
			Confirmations: 12,
			// FIXME
			ExternalSystem: 9,
		}

		err := figure.
			Out(&config).
			With(figure.BaseHooks, utils.ETHHooks).
			From(app.Config(ctx).Get(conf.ServiceERC20Deposit)).
			Please()
		if err != nil {
			return nil, errors.Wrap(err, "failed to figure out")
		}

		horizon := app.Config(ctx).Horizon().WithSigner(config.Signer)

		addrProvider := addrstate.New(
			ctx,
			app.Log(ctx),
			[]addrstate.StateMutator{
				addrstate.ExternalSystemBindingMutator(config.ExternalSystem),
				addrstate.BalanceMutator(config.DepositAsset),
			},
			horizon.Listener(),
		)

		builder, err := horizon.TXBuilder()
		if err != nil {
			return nil, errors.Wrap(err, "failed to init tx builder")
		}

		eth := app.Config(ctx).Ethereum()

		return deposit.New(&deposit.Opts{
			app.Log(ctx),
			config.Source,
			config.Signer,
			conf.ServiceERC20Deposit,
			conf.ServiceERC20DepositVerify,
			config.Cursor,
			config.Confirmations,
			app.Config(ctx).Horizon().WithSigner(config.Signer),
			config.ExternalSystem,
			addrProvider,
			app.Config(ctx).Discovery(),
			builder,
			internal.NewERC20Helper(eth, config.DepositAsset, config.Token),
			config.DisableVerify,
		}), nil
	})
}
