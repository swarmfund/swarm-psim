package conf

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
)

func (c *ViperConfig) Bitcoin() *bitcoin.Client {
	c.Lock()
	defer c.Unlock()

	if c.btcClient != nil {
		return c.btcClient
	}

	config := bitcoin.ConnectorConfig{}

	err := figure.Out(&config).From(c.GetRequired("bitcoin")).With(figure.BaseHooks, bitcoin.FigureHooks).Please()
	if err != nil {
		panic(errors.Wrap(err, "Failed to parse bitcoin config entry"))
	}

	connector := bitcoin.NewNodeConnector(config)
	c.btcClient = bitcoin.NewClient(connector)

	return c.btcClient
}
