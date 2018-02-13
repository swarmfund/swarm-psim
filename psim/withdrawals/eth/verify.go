package eth

import (
	"context"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/swarmfund/psim/psim/withdrawals/eth/internal"
	"gitlab.com/swarmfund/psim/psim/withdveri"
)

func init() {
	app.RegisterService(conf.ServiceETHWithdrawVerify, func(ctx context.Context) (app.Service, error) {
		var config VerifyConfig
		err := figure.
			Out(&config).
			From(app.Config(ctx).Get(conf.ServiceETHWithdrawVerify)).
			With(figure.BaseHooks, utils.ETHHooks).
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
			return nil, errors.Wrap(err, "Failed to get Horizon info")
		}

		builder := xdrbuild.NewBuilder(info.Passphrase, info.TXExpirationPeriod)

		wallet := eth.NewWallet()
		address, err := wallet.ImportHEX(config.Key)
		if err != nil {
			return nil, errors.Wrap(err, "failed to import key")
		}

		token, err := internal.NewToken(*config.Token, internal.ERC20ABI)
		if err != nil {
			return nil, errors.Wrap(err, "failed to init token")
		}

		ethClient := app.Config(ctx).Ethereum()

		return withdveri.New(
			conf.ServiceETHWithdrawVerify,
			app.Log(ctx),
			config.Source,
			config.Signer,
			horizon,
			builder,
			listener,
			app.Config(ctx).Discovery(),
			internal.NewHelper(config.Asset, config.Threshold, ethClient, address, wallet, config.GasPrice, token),
		), nil
	})
}
