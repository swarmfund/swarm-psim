package ethsupervisor

import (
	"context"

	"fmt"

	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/psim/addrstate"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/ethsupervisor/internal"
	"gitlab.com/swarmfund/psim/psim/supervisor"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceETHSupervisor, func(ctx context.Context) (app.Service, error) {
		config := Config{
			Supervisor: supervisor.NewConfig(conf.ServiceETHSupervisor),
		}

		err := figure.
			Out(&config).
			From(app.Config(ctx).GetRequired(conf.ServiceETHSupervisor)).
			With(supervisor.DLFigureHooks, figure.BaseHooks, utils.ETHHooks).
			Please()
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to figure out %s", conf.ServiceETHSupervisor))
		}

		if config.FixedDepositFee == nil {
			return nil, errors.New("'fixed_deposit_fee' cannot be empty")
		}

		commonSupervisor, err := supervisor.InitNew(ctx, conf.ServiceETHSupervisor, config.Supervisor)
		if err != nil {
			return nil, errors.Wrap(err, "failed to init supervisor common")
		}

		ethClient := app.Config(ctx).Ethereum()

		horizon := app.Config(ctx).Horizon().WithSigner(config.Supervisor.SignerKP)

		log := app.Log(ctx)
		state := addrstate.New(
			ctx,
			log.WithField("service", "addrstate"),
			internal.StateMutator(config.BaseAsset, config.DepositAsset),
			horizon.Listener(),
		)

		return New(commonSupervisor, ethClient, state, config), nil
	})
}

type Service struct {
	*supervisor.Service
	eth    *ethclient.Client
	state  State
	config Config

	// internal state
	txCh     chan internal.Transaction
	blocksCh chan uint64

	// config
	depositThreshold *big.Int
}

// FIXME Hardcoded threshold
func New(supervisor *supervisor.Service, eth *ethclient.Client, state State, config Config) *Service {
	s := &Service{
		Service: supervisor,
		eth:     eth,
		state:   state,
		config:  config,
		// could be buffered to increase throughput
		txCh:     make(chan internal.Transaction, 1),
		blocksCh: make(chan uint64, 1),
		// FIXME Hardcoded threshold
		depositThreshold: big.NewInt(1000000000000),
	}

	s.AddRunner(s.watchHeight)
	for i := 0; i < 10; i++ {
		s.AddRunner(s.processBlocks)
	}
	s.AddRunner(s.processTXs)

	return s
}
