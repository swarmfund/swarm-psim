package ethsupervisor

import (
	"context"

	"fmt"

	"encoding/json"
	"net/http"

	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"gitlab.com/tokend/psim/addrstate"
	"gitlab.com/tokend/psim/figure"
	"gitlab.com/tokend/psim/psim/app"
	"gitlab.com/tokend/psim/psim/conf"
	"gitlab.com/tokend/psim/psim/ethsupervisor/internal"
	"gitlab.com/tokend/psim/psim/supervisor"
	"gitlab.com/tokend/psim/psim/utils"
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

		horizonRequest := func(method, endpoint string, target interface{}) error {
			r, err := horizon.SignedRequest(method, endpoint, config.Supervisor.SignerKP)
			if err != nil {
				return errors.Wrap(err, "failed to init request")
			}
			response, err := http.DefaultClient.Do(r)
			if err != nil {
				return errors.Wrap(err, "request failed")
			}
			defer response.Body.Close()
			if err := json.NewDecoder(response.Body).Decode(&target); err != nil {
				return errors.Wrap(err, "failed to unmarshal")
			}
			return nil
		}

		state := addrstate.New(
			internal.StateMutator,
			addrstate.NewLedgersProvider(app.Log(ctx).WithField("service", "eth-ledger-provider"), horizonRequest),
			addrstate.NewChangesProvider(app.Log(ctx).WithField("service", "eth-changes-provider"), horizonRequest),
		)

		return New(commonSupervisor, ethClient, state), nil
	})
}

type Service struct {
	*supervisor.Service
	eth   *ethclient.Client
	state State

	// internal state
	txCh     chan internal.Transaction
	blocksCh chan uint64
	// config
	depositTreshold big.Int
}

func New(supervisor *supervisor.Service, eth *ethclient.Client, state State) *Service {
	s := &Service{
		Service:  supervisor,
		eth:      eth,
		state:    state,
		txCh:     make(chan internal.Transaction),
		blocksCh: make(chan uint64),
	}

	s.AddRunner(s.watchHeight)
	for i := 0; i < 10; i++ {
		s.AddRunner(s.processBlocks)
	}
	s.AddRunner(s.processTXs)

	return s
}
