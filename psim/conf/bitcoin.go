package conf

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/figure"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
)

var (
	btcClient *bitcoin.Client
)

func (c *ViperConfig) Bitcoin() *bitcoin.Client {
	if btcClient == nil {
		config := bitcoin.ConnectorConfig{}

		err := figure.Out(&config).From(c.GetRequired("bitcoin")).With(figure.BaseHooks).Please()
		if err != nil {
			panic(errors.Wrap(err, "Failed to parse bitcoin config entry"))
		}

		connector := bitcoin.NewNodeConnector(config)
		btcClient = bitcoin.NewClient(connector)
	}

	return btcClient
}
