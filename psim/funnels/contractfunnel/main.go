package contractfunnel

import (
	"context"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/addrstate"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
)

func init() {
	app.RegisterService(conf.ServiceETHContractFunnel, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	config, err := NewConfig(app.Config(ctx).Get(conf.ServiceETHContractFunnel))
	if err != nil {
		return nil, errors.Wrap(err, "failed to init config")
	}

	keypair, err := eth.NewKeypair(config.PrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init keypair")
	}

	horizon := app.Config(ctx).Horizon().WithSigner(config.Signer)

	mutators := make([]addrstate.StateMutator, 0, len(config.ExternalSystems))
	for _, system := range config.ExternalSystems {
		mutators = append(mutators, addrstate.ExternalSystemBindingMutator{SystemType: system})
	}

	addrProvider := addrstate.New(
		ctx,
		app.Log(ctx),
		mutators,
		horizon.Listener(),
	)

	return &Service{
		Opts: Opts{
			app.Log(ctx),
			config,
			app.Config(ctx).Ethereum(),
			keypair,
			addrProvider,
			config.ExternalSystems,
			config.Tokens,
			config.HotWallet,
			config.Threshold,
			config.GasPrice,
		},
	}, nil
}
