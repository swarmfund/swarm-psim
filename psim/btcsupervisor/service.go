package btcsupervisor

import (
	"context"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/supervisor"
)

// BTCClient must be implemented by a BTC Client to pass into Service constructor.
type BTCClient interface {
	GetBlockCount() (uint64, error)
	GetBlock(blockIndex uint64) (*btcutil.Block, error)
	GetNetParams() *chaincfg.Params
}

// AddressProvider must be implemented by WatchAddress storage to pass into Service constructor.
type AddressProvider interface {
	AddressAt(ctx context.Context, t time.Time, btcAddress string) (tokendAddress *string)
}

// Service implements app.Service interface, it supervises Stripe transactions
// and send CoinEmissionRequests to Horizon if arrived Charge detected.
//
// Service uses supervisor.Service for common for supervisors logic, such as Leadership and Profiling.
type Service struct {
	*supervisor.Service

	// TODO Interface?
	horizon         *horizon.Connector
	config          Config
	btcClient       BTCClient
	addressProvider AddressProvider
}

// New is constructor for the btcsupervisor Service.
func New(commonSupervisor *supervisor.Service, config Config, btcClient BTCClient, addressProvider AddressProvider, horizon *horizon.Connector) *Service {
	result := &Service{
		Service: commonSupervisor,

		horizon:         horizon,
		config:          config,
		btcClient:       btcClient,
		addressProvider: addressProvider,
	}

	result.AddRunner(result.processBTCBlocksInfinitely)

	return result
}
