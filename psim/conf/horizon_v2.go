package conf

import (
	"net/url"

	horizon "gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var (
	horizonV2 *horizon.Connector
)

func (c *ViperConfig) HorizonV2() *horizon.Connector {
	if horizonV2 != nil {
		return horizonV2
	}

	v := c.viper.Sub("horizon")
	if v == nil {
		panic("config entry is missing")
	}

	addr := v.GetString("addr")
	endpoint, err := url.Parse(addr)
	if err != nil {
		panic(errors.Wrap(err, "Failed to parse addr into url", logan.F{
			"addr": addr,
		}))
	}

	horizonV2 = horizon.NewConnector(endpoint)

	return horizonV2
}
