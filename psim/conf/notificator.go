package conf

import (
	"net/url"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/notificator-server/client"
	"gitlab.com/swarmfund/psim/psim/utils"
)

type NotificatorConfig struct {
	URL    *url.URL `fig:"url,required"`
	Secret string   `fig:"secret"`
	Public string   `fig:"public"`
}

func (c *ViperConfig) Notificator() *notificator.Connector {
	c.Lock()
	defer c.Unlock()

	if c.notificatorClient != nil {
		return c.notificatorClient
	}
	conf := NotificatorConfig{}

	err := figure.
		Out(&conf).
		From(c.Get("notificator")).
		With(figure.BaseHooks, utils.CommonHooks).
		Please()

	if err != nil {
		panic(errors.Wrap(err, "Failed to figure out Notificator"))
	}

	client := notificator.NewConnector(notificator.Pair{
		Secret: conf.Secret,
		Public: conf.Public,
	}, *conf.URL)

	c.notificatorClient = client

	return c.notificatorClient
}
