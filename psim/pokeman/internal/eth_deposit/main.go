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

		esProvider := NewExternalSystemProvider(horizon.Assets(), config.Asset)
		esType, err := esProvider.GetExternalSystemType()

		currentBalanceProvider := NewCurrentBalanceProvider(horizon, config.Source.Address(), config.Asset)
		b, err := currentBalanceProvider.CurrentBalance()

		return &Service{
			app.Log(ctx),
			NewBalancePoller(ctx, app.Log(ctx), 30, currentBalanceProvider),
			NewEthTxProvider(app.Config(ctx).Ethereum(), ctx, config.Keypair, app.Log(ctx)),
			NewNativeTxProvider(horizon, builder, config.Source, config.Keypair, config.Signer, config.Asset, b.BalanceID, ctx),
			NewExternalBindingDataProvider(horizon, config.Source.Address(), esType),
			NewExternalSystemBinder(builder, horizon, config.Source, config.Signer, esType),
			currentBalanceProvider,
			esProvider,
		}, nil
	})
}
