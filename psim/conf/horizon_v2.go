package conf

import (
	"net/url"

	"github.com/pkg/errors"
	horizon "gitlab.com/swarmfund/horizon-connector/v2"
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

	endpoint, err := url.Parse(v.GetString("addr"))
	if err != nil {
		panic(errors.Wrap(err, "failed to parse addr"))
	}

	horizonV2 = horizon.NewConnector(endpoint)

	return horizonV2
}
