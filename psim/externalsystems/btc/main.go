package btc

import (
	"context"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/externalsystems/deployer"
	"gitlab.com/swarmfund/psim/psim/externalsystems/derive"
	"gitlab.com/swarmfund/psim/psim/internal"
)

func init() {
	app.RegisterService(conf.ServiceBTCDeployer, func(ctx context.Context) (app.Service, error) {
		config, err := NewConfig(app.Config(ctx).Get(conf.ServiceBTCDeployer))
		if err != nil {
			return nil, errors.Wrap(err, "failed to init config")
		}

		horizon := app.Config(ctx).Horizon().WithSigner(config.Signer)

		builder, err := horizon.TXBuilder()
		if err != nil {
			return nil, errors.Wrap(err, "failed init tx builder")
		}

		deployerID := internal.Hash64(base58.Decode(config.HDKey))

		deriver, err := derive.NewBTCFamilyDeriver(derive.NetworkTypeBTCMainnet, config.HDKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to init deriver")
		}
		return deployer.NewService(deployer.Opts{
			Log:           app.Log(ctx),
			ExternalTypes: config.ExternalTypes,
			EntityCount:   deployer.ExternalSystemPoolEntityCount(horizon),
			TargetCount:   config.TargetCount,
			Deriver:       deriver,
			TXBuilder:     builder,
			Source:        config.Source,
			Signer:        config.Signer,
			Horizon:       horizon,
			DeployerID:    deployerID,
		}), nil
	})
}
