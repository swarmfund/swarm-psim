package conf

import (
	"gitlab.com/tokend/psim/psim/bitcoin"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/psim/figure"
)

var (
	btcClient *bitcoin.Client
)

func (c *ViperConfig) Bitcoin() (*bitcoin.Client, error) {
	if btcClient == nil {
		config := bitcoin.ConnectorConfig{}
		err := figure.Out(&config).From(c.Get("bitcoin")).With(figure.BaseHooks).Please()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to parse bitcoin config entry")
		}

		connector := bitcoin.NewNodeConnector(config)
		btcClient = bitcoin.NewClient(connector)
	}

	return btcClient, nil
}
