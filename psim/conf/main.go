package conf

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/client"
	discovery "gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/notificator-server/client"
	"gitlab.com/swarmfund/psim/mixpanel"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
	"gitlab.com/swarmfund/psim/psim/notifications"
	"gitlab.com/swarmfund/psim/salesforce"
	"gitlab.com/tokend/horizon-connector"
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
	// TODO Consider creating HorizonWithSigner() method.
	Horizon() *horizon.Connector
	Services() []string
	// TODO Panic instead of returning errors.
	Stripe() (*client.API, error)
	Ethereum() *ethclient.Client
	Bitcoin() *bitcoin.Client
	Notificator() *notificator.Connector
	NotificationSender() *notifications.SlackSender
	S3() *session.Session
	Mixpanel() *mixpanel.Connector // TODO
	Salesforce() *salesforce.Connector
}

type ViperConfig struct {
	viper *viper.Viper

	// internal singletons
	*sync.Mutex
	horizon            *horizon.Connector
	session            *session.Session
	btcClient          *bitcoin.Client
	discoveryClient    *discovery.Client
	defaultLog         *logan.Entry
	notificationSender *notifications.SlackSender
	notificatorClient  *notificator.Connector
	stripeClient       *client.API
	salesforce         *salesforce.Connector
	mixpanel           *mixpanel.Connector
}

func NewViperConfig(fn string) *ViperConfig {
	config := ViperConfig{
		viper: viper.GetViper(),
		Mutex: &sync.Mutex{},
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
