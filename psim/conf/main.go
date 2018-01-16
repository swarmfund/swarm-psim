package conf

import (
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/client"
	discovery "gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/notificator-server/client"
	horizon "gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
)

// TODO: viper's Get* methods won't throw error if value is invalid

type Config interface {
	// TODO Panic instead of returning errors.
	Init() error
	// DEPRECATED Use GetRequired instead (it panics if key is missing).
	Get(key string) map[string]interface{}
	GetRequired(key string) map[string]interface{}
	Discovery() *discovery.Client
	// TODO Panic instead of returning errors.
	Log() (*logan.Entry, error)
	Horizon() *horizon.Connector
	Services() []string
	Stripe() (*client.API, error)
	Ethereum() *ethclient.Client
	Bitcoin() *bitcoin.Client
	Notificator() (*notificator.Connector, error)
}

type ViperConfig struct {
	viper *viper.Viper

	// internal singletons
	*sync.Mutex
	horizon *horizon.Connector
}

func NewViperConfig(fn string) *ViperConfig {
	config := ViperConfig{
		viper: viper.GetViper(),
	}
	config.viper.SetConfigFile(fn)
	return &config
}

func (c *ViperConfig) Init() error {
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "failed to set config file")
	}
	return nil
}

// GetRequired panics, if key is missing.
func (c *ViperConfig) GetRequired(key string) map[string]interface{} {
	v := c.viper.Sub(key)
	if v == nil {
		panic(errors.From(errors.New("Config entry is missing."), logan.F{
			"config_key": key,
		}))
	}

	return c.viper.GetStringMap(key)
}

func (c *ViperConfig) Get(key string) map[string]interface{} {
	return c.viper.GetStringMap(key)
}
