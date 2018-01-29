package btcsupervisor

import (
	"context"

	"fmt"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/addrstate"
	"gitlab.com/swarmfund/psim/figure"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/btcsupervisor/internal"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/supervisor"
)

func init() {
	app.RegisterService(conf.ServiceBTCSupervisor, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)

	config := Config{
		Supervisor: supervisor.NewConfig(conf.ServiceBTCSupervisor),
	}

	err := figure.
		Out(&config).
		From(globalConfig.GetRequired(conf.ServiceBTCSupervisor)).
		With(supervisor.ConfigFigureHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to figure out %s", conf.ServiceBTCSupervisor))
	}

	commonSupervisor, err := supervisor.InitNew(ctx, conf.ServiceBTCSupervisor, config.Supervisor)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to init common Supervisor")
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.Supervisor.SignerKP)

	log := app.Log(ctx)

	addressProvider := addrstate.New(
		ctx,
		log.WithField("service", "addrstate"),
		internal.StateMutator,
		horizonConnector.Listener(),
	)

	return New(commonSupervisor, config, globalConfig.Bitcoin(), addressProvider, horizonConnector), nil
}
