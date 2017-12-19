package ethsupervisor

import (
	"context"

	"fmt"

	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	horizon "gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/psim/addrstate"
	"gitlab.com/swarmfund/psim/figure"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/ethsupervisor/internal"
	"gitlab.com/swarmfund/psim/psim/supervisor"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceETHSupervisor, func(ctx context.Context) (utils.Service, error) {
		config := Config{
			Supervisor: supervisor.NewConfig(conf.ServiceETHSupervisor),
		}

		err := figure.
			Out(&config).
			From(app.Config(ctx).Get(conf.ServiceETHSupervisor)).
			With(supervisor.ConfigFigureHooks).
			Please()
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to figure out %s", conf.ServiceETHSupervisor))
		}

		commonSupervisor, err := supervisor.InitNew(ctx, conf.ServiceETHSupervisor, config.Supervisor)
		if err != nil {
			return nil, errors.Wrap(err, "failed to init supervisor common")
		}

		ethClient := app.Config(ctx).Ethereum()

		horizon, err := app.Config(ctx).Horizon()
		if err != nil {
			return nil, errors.Wrap(err, "failed to init horizon")
		}

		horizonV2 := app.Config(ctx).HorizonV2()

		log := app.Log(ctx)
		state := addrstate.New(
			ctx,
			log.WithField("service", "addrstate"),
			internal.StateMutator,
			horizonV2.Listener(),
		)

		return New(commonSupervisor, ethClient, state, config, horizon), nil
	})
}

type Service struct {
	*supervisor.Service
	eth     *ethclient.Client
	state   State
	config  Config
	horizon *horizon.Connector

	// internal state
	txCh     chan internal.Transaction
	blocksCh chan uint64

	// config
	depositThreshold *big.Int
}

func New(supervisor *supervisor.Service, eth *ethclient.Client, state State, config Config, horizon *horizon.Connector) *Service {
	s := &Service{
		Service: supervisor,
		eth:     eth,
		state:   state,
		config:  config,
		horizon: horizon,
		// could be buffered to increase throughput
		txCh:     make(chan internal.Transaction),
		blocksCh: make(chan uint64),
		// FIXME
		depositThreshold: big.NewInt(1000000000000),
	}

	s.AddRunner(s.watchHeight)
	for i := 0; i < 10; i++ {
		s.AddRunner(s.processBlocks)
	}
	s.AddRunner(s.processTXs)

	return s
}
