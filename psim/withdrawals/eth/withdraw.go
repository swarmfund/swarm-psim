package eth

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/swarmfund/psim/psim/withdrawals/withdraw"
	"gitlab.com/swarmfund/psim/psim/withdrawals/eth/internal"
)

func init() {
	app.RegisterService(conf.ServiceETHWithdraw, func(ctx context.Context) (app.Service, error) {
		var config WithdrawConfig
		err := figure.
			Out(&config).
			With(figure.BaseHooks, utils.ETHHooks).
			From(app.Config(ctx).Get(conf.ServiceETHWithdraw)).
			Please()

		horizon := app.Config(ctx).Horizon().WithSigner(config.Signer)

		info, err := horizon.Info()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get horizon info")
		}
		builder := xdrbuild.NewBuilder(info.Passphrase, info.TXExpirationPeriod)

		ethClient := app.Config(ctx).Ethereum()

		wallet := eth.NewWallet()
		address, err := wallet.ImportHEX(config.Key)
		if err != nil {
			return nil, errors.Wrap(err, "failed to import key")
		}

		var token *internal.Token
		if config.Token != nil {
			token, err = internal.NewToken(*config.Token, internal.ERC20ABI)
			if err != nil {
				return nil, errors.Wrap(err, "failed to init token")
			}
		}

		return withdraw.New(
			conf.ServiceETHWithdraw,
			conf.ServiceETHWithdrawVerify,
			config.Signer,
			app.Log(ctx),
			horizon.Listener(),
			horizon,
			builder,
			app.Config(ctx).Discovery(),
			internal.NewHelper(config.Asset, config.Threshold, ethClient, address, wallet, config.GasPrice, token, app.Log(ctx)),
		), nil
	})
}
