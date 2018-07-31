package eth_deposit

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/internal"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/horizon-connector"
	)

func init() {
	app.RegisterService(conf.PokemanETHDepositService, func(ctx context.Context) (app.Service, error) {
		config, err := NewConfig(app.Config(ctx).Get(conf.PokemanETHDepositService))
		if err != nil {
			return nil, errors.Wrap(err, "failed to init config")
		}

		connector := app.Config(ctx).Horizon()

		builder, err := connector.TXBuilder()
		if err != nil {
			return nil, errors.Wrap(err, "failed to init tx builder")
		}

		log := app.Log(ctx)

		return (&Service{
			log,
			func() (int32, error) {
				return internal.GetExternalSystemType(connector.Assets(), config.Asset)
			},
			func() (horizon.Balance, error) {
				return connector.Accounts().CurrentBalanceIn(config.Source.Address(), config.Asset)
			},
			func(externalSystem int32) (*string, error) {
			return connector.Accounts().CurrentExternalBindingData(config.Source.Address(), externalSystem)
			},
			connector.Submitter().Submit,
			func(op xdrbuild.Operation) (string, error) {
				return builder.Transaction(config.Signer).Op(op).Sign(config.Signer).Marshal()
			},
			NewNativeTxProvider(connector, builder, config.Source, config.Keypair, config.Signer, config.Asset, config.Source.Address()),
			NewEthTxProvider(app.Config(ctx).Ethereum(), config.Keypair, log),
		}).WithTimeout(config.PollingTimeout), nil
	})
}
