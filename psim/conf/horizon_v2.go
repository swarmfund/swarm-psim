package conf

import (
	"net/url"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	horizon "gitlab.com/tokend/horizon-connector"
)

func (c *ViperConfig) Horizon() *horizon.Connector {
	c.Lock()
	defer c.Unlock()

	if c.horizon == nil {
		v := c.viper.Sub("horizon")
		if v == nil {
			panic("horizon config entry is missing")
		}

		addr := v.GetString("addr")
		endpoint, err := url.Parse(addr)
		if err != nil {
			panic(errors.Wrap(err, "failed to parse addr", logan.F{
				"addr": addr,
			}))
		}

		c.horizon = horizon.NewConnector(endpoint)
	}

	return c.horizon
}
