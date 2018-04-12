package erc20

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/deposits/erc20/internal"
	"gitlab.com/swarmfund/psim/psim/deposits/depositveri"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceERC20DepositVerify, func(ctx context.Context) (app.Service, error) {
		var config VerifyConfig
		err := figure.
			Out(&config).
			With(figure.BaseHooks, utils.ETHHooks).
			From(app.Config(ctx).Get(conf.ServiceERC20DepositVerify)).
			Please()
		if err != nil {
			return nil, errors.Wrap(err, "failed to figure out")
		}

		listener, err := ape.Listener(config.Host, config.Port)
		if err != nil {
			return nil, errors.Wrap(err, "failed to init listener")
		}

		horizon := app.Config(ctx).Horizon().WithSigner(config.Signer)

		info, err := horizon.Info()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get horizon info")
		}

		builder := xdrbuild.NewBuilder(info.Passphrase, info.TXExpirationPeriod)
		eth := app.Config(ctx).Ethereum()

		return depositveri.New(
			"ethereum",
			conf.ServiceERC20DepositVerify,
			app.Log(ctx),
			config.Signer,
			config.Confirmations,
			horizon,
			builder,
			listener,
			app.Config(ctx).Discovery(),
			internal.NewERC20Helper(eth, config.DepositAsset, config.Token),
		), nil
	})
}
