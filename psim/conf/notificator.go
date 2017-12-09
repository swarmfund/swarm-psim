package conf

import (
	"net/url"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/notificator-server/client"
	"gitlab.com/tokend/psim/figure"
	"gitlab.com/tokend/psim/psim/utils"
)

var (
	notificatorClient *notificator.Connector
)

type NotificatorConfig struct {
	URL    *url.URL
	Secret string
	Public string
}

func (c *ViperConfig) Notificator() (*notificator.Connector, error) {
	if notificatorClient != nil {
		return notificatorClient, nil
	}
	conf := NotificatorConfig{}

	err := figure.
		Out(&conf).
		From(c.Get("notificator")).
		With(figure.BaseHooks, utils.CommonHooks).
		Please()

	if err != nil {
		return nil, errors.Wrap(err, "failed to figure out notificator")
	}

	client := notificator.NewConnector(notificator.Pair{
		Secret: conf.Secret,
		Public: conf.Public,
	}, *conf.URL)

	notificatorClient = client

	return notificatorClient, nil
}
