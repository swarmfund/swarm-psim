package eth_deposit

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
)

func init() {
	app.RegisterService(conf.PokemanETHDepositService, func(ctx context.Context) (app.Service, error) {
		config, err := NewConfig(app.Config(ctx).Get(conf.PokemanETHDepositService))
		if err != nil {
			return nil, errors.Wrap(err, "failed to init config")
		}

		horizon := app.Config(ctx).Horizon()

		builder, err := horizon.TXBuilder()
		if err != nil {
			return nil, errors.Wrap(err, "failed to init tx builder")
		}

		log := app.Log(ctx)

		ethTxProvider := NewEthTxProvider(app.Config(ctx).Ethereum(), config.Keypair, log)

		return NewService(
			log,
			ethTxProvider,
			app.Config(ctx).Slack(),
			horizon,
			config,
			builder,
		), nil
	})
}
