package ethcontracts

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/internal"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
	"gitlab.com/tokend/horizon-connector"
)

func init() {
	app.RegisterService(conf.ServiceETHContracts, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	config, err := NewConfig(app.Config(ctx).Get(conf.ServiceETHContracts))
	if err != nil {
		return nil, errors.Wrap(err, "failed to init config")
	}

	keypair, err := eth.NewKeypair(config.ETHPrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init keypair")
	}

	horizon := app.Config(ctx).Horizon().WithSigner(config.Signer)

	builder, err := horizon.TXBuilder()
	if err != nil {
		return nil, errors.Wrap(err, "failed init tx builder")
	}

	deployerID := internal.Hash64(keypair.Address().Bytes())

	return NewService(
		app.Log(ctx),
		config,
		//globalConfig.Ethereum(),
		builder,
		horizon,
		keypair,
		deployerID,
		ExternalSystemPoolEntityCount(horizon),
		app.Config(ctx).Ethereum(),
	), nil
}

func ExternalSystemPoolEntityCount(horizon *horizon.Connector) func(string) (uint64, error) {
	return func(systemType string) (uint64, error) {
		stats, err := horizon.System().Statistics()
		if err != nil {
			return 0, errors.Wrap(err, "failed to get system stats")
		}
		count := stats.ExternalSystemPoolEntriesCount[systemType]
		return count, nil
	}
}
