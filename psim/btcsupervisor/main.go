package btcsupervisor

import (
	"context"

	"fmt"

	"time"

	"github.com/piotrnar/gocoin/lib/btc"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/figure"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/btcsupervisor/internal"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/create_account_streamer"
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
			return nil, errors.Wrap(err, fmt.Sprintf("failed to figure out %s", conf.ServiceBTCSupervisor))
		}

		commonSupervisor, err := supervisor.InitNew(ctx, conf.ServiceBTCSupervisor, config.Supervisor)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to init common supervisor")
		}

		btcClient, err := globalConfig.Bitcoin()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get Bitcoin client")
		}

		horizonConnector, err := app.Config(ctx).Horizon()
		if err != nil {
			panic(err)
		}

		createAccountStreamer := create_account_streamer.New(app.Log(ctx), config.Supervisor.SignerKP, horizonConnector,
			5*time.Second)
		addressQ := internal.NewAddressQ(ctx, createAccountStreamer)

		return New(commonSupervisor, config, btcClient, addressQ), nil
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
type AddressQ interface {
	// If this btc Address - "" must be returned.
	GetAccountID(btcAddress string) string
	ReadinessWaiter() <-chan struct{}
	Run()
}

// Service implements utils.Service interface, it supervises Stripe transactions
// and send CoinEmissionRequests to Horizon if arrived Charge detected.
//
// Service uses supervisor.Service for common for supervisors logic, such as Leadership and Profiling.
type Service struct {
	*supervisor.Service

	config    Config
	btcClient BTCClient
	addressQ  AddressQ
}

// New is constructor for the btcsupervisor Service.
func New(commonSupervisor *supervisor.Service, config Config, btcClient BTCClient, addressQ AddressQ) *Service {
	result := &Service{
		Service: commonSupervisor,

		config:    config,
		btcClient: btcClient,
		addressQ:  addressQ,
	}

	result.AddRunner(result.addressQ.Run)
	result.AddRunner(result.processBTCBlocksInfinitely)

	return result
}
