package btcsupervisor

import (
	"gitlab.com/swarmfund/psim/psim/supervisor"
	"gitlab.com/swarmfund/horizon-connector"
	"time"
	"context"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

// BTCClient must be implemented by a BTC Client to pass into Service constructor.
type BTCClient interface {
	GetBlockCount() (uint64, error)
	GetBlock(blockIndex uint64) (*btcutil.Block, error)
	GetNetParams() *chaincfg.Params
}

// AddressQ must be implemented by WatchAddress storage to pass into Service constructor.
type AccountDataProvider interface {
	AddressAt(ctx context.Context, t time.Time, btcAddress string) (tokendAddress *string)
	PriceAt(ctx context.Context, ts time.Time) *int64
}

// Service implements app.Service interface, it supervises Stripe transactions
// and send CoinEmissionRequests to Horizon if arrived Charge detected.
//
// Service uses supervisor.Service for common for supervisors logic, such as Leadership and Profiling.
type Service struct {
	*supervisor.Service

	// TODO Interface?
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
