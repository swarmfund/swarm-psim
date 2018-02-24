package conf

import (
	"net/url"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/notificator-server/client"
	"gitlab.com/swarmfund/psim/figure"
	"gitlab.com/swarmfund/psim/psim/utils"
)

var (
	notificatorClient *notificator.Connector
)

type NotificatorConfig struct {
	URL    *url.URL
	Secret string
	Public string
}

func (c *ViperConfig) Notificator() *notificator.Connector {
	if notificatorClient != nil {
		return notificatorClient
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

	notificatorClient = client

	return notificatorClient
}
