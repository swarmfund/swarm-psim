package btcsupervisor

import (
	"context"

	"fmt"

	"time"

	"github.com/piotrnar/gocoin/lib/btc"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/psim/addrstate"
	"gitlab.com/swarmfund/psim/figure"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/btcsupervisor/internal"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/horizonreq"
	"gitlab.com/swarmfund/psim/psim/supervisor"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	setupFn := func(ctx context.Context) (utils.Service, error) {
		globalConfig := app.Config(ctx)

		config := Config{
			Supervisor: supervisor.NewConfig(conf.ServiceBTCSupervisor),
		}

		err := figure.
			Out(&config).
			From(globalConfig.Get(conf.ServiceBTCSupervisor)).
			With(supervisor.ConfigFigureHooks).
			Please()
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Failed to figure out %s", conf.ServiceBTCSupervisor))
		}

		commonSupervisor, err := supervisor.InitNew(ctx, conf.ServiceBTCSupervisor, config.Supervisor)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to init common Supervisor")
		}

		btcClient, err := globalConfig.Bitcoin()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get Bitcoin client from global config")
		}

		horizonConnector, err := app.Config(ctx).Horizon()
		if err != nil {
			panic(err)
		}

		log := app.Log(ctx)
		requester := horizonreq.NewHorizonRequester(horizonConnector, config.Supervisor.SignerKP)
		addressProvider := addrstate.New(log, internal.StateMutator,
			addrstate.NewLedgersProvider(log.WithField("service", "btc-ledger-provider"), requester),
			addrstate.NewChangesProvider(log.WithField("service", "btc-changes-provider"), requester),
			requester,
		)

		return New(commonSupervisor, config, btcClient, addressProvider, horizonConnector), nil
	}

	app.RegisterService(conf.ServiceBTCSupervisor, setupFn)
}

// BTCClient must be implemented by a BTC Client to pass into Service constructor.
type BTCClient interface {
	IsTestnet() bool
	GetBlockCount() (uint64, error)
	GetBlock(blockIndex uint64) (*btc.Block, error)
}

// AddressQ must be implemented by WatchAddress storage to pass into Service constructor.
type AccountDataProvider interface {
	AddressAt(ctx context.Context, t time.Time, btcAddress string) (tokendAddress *string)
	PriceAt(ctx context.Context, ts time.Time) *int64
}

// Service implements utils.Service interface, it supervises Stripe transactions
// and send CoinEmissionRequests to Horizon if arrived Charge detected.
//
// Service uses supervisor.Service for common for supervisors logic, such as Leadership and Profiling.
type Service struct {
	*supervisor.Service

	horizon             *horizon.Connector
	config              Config
	btcClient           BTCClient
	accountDataProvider AccountDataProvider
}

// New is constructor for the btcsupervisor Service.
func New(commonSupervisor *supervisor.Service, config Config, btcClient BTCClient, addressProvider AccountDataProvider, horizon *horizon.Connector) *Service {
	result := &Service{
		Service: commonSupervisor,

		horizon:             horizon,
		config:              config,
		btcClient:           btcClient,
		accountDataProvider: addressProvider,
	}

	result.AddRunner(result.processBTCBlocksInfinitely)

	return result
}
