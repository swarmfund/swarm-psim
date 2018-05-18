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
	"gitlab.com/tokend/go/xdrbuild"
)

func init() {
	app.RegisterService(conf.ServiceERC20Deposit, func(ctx context.Context) (app.Service, error) {
		config := DepositConfig{
			Confirmations: 12,
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
			internal.StateMutator(config.BaseAsset, config.DepositAsset),
			horizon.Listener(),
		)

		info, err := horizon.Info()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get horizon info")
		}
		builder := xdrbuild.NewBuilder(info.Passphrase, info.TXExpirationPeriod)

		ethclient := app.Config(ctx).Ethereum()

		return deposit.New(&deposit.Opts{
			app.Log(ctx),
			config.Source,
			config.Signer,
			conf.ServiceERC20Deposit,
			conf.ServiceERC20DepositVerify,
			config.Cursor,
			config.Confirmations,
			app.Config(ctx).Horizon().WithSigner(config.Signer),
			addrProvider,
			app.Config(ctx).Discovery(),
			builder,
			internal.NewERC20Helper(ethclient, config.DepositAsset, config.Token),
		}), nil
	})
}
