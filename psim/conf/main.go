package conf

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/client"
	discovery "gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/notificator-server/client"
	horizon "gitlab.com/swarmfund/horizon-connector"
	horizonv2 "gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
)

// TODO: viper's Get* methods won't throw error if value is invalid

type Config interface {
	Init() error
	Get(key string) map[string]interface{}
	Discovery() (*discovery.Client, error)
	Log() (*logan.Entry, error)
	Horizon() (*horizon.Connector, error)
	HorizonV2() *horizonv2.Connector
	Services() []string
	Stripe() (*client.API, error)
	Ethereum() *ethclient.Client
	Bitcoin() (*bitcoin.Client, error)
	Notificator() (*notificator.Connector, error)
}

type ViperConfig struct {
	viper *viper.Viper
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

func (c *ViperConfig) Get(key string) map[string]interface{} {
	return c.viper.GetStringMap(key)
}
