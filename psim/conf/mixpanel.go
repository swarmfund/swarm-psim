package conf

import (
	"net/url"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/mixpanel"
)

// TODO use figure out

// Salesforce returns a ready-to-use salesforce connector
func (c *ViperConfig) Mixpanel() *mixpanel.Connector {
	c.Lock()
	defer c.Unlock()

	if c.mixpanel == nil {
		v := c.viper.Sub("mixpanel")
		if v == nil {
			panic("mixpanel config entry is missing")
		}

		apiRawURL := v.GetString("api_url")
		apiURL, err := url.Parse(apiRawURL)
		if err != nil {
			panic(errors.Wrap(err, "failed to parse mixpanel api url", logan.F{
				"api_url": apiRawURL,
			}))
		}
		token := v.GetString("token")

		c.mixpanel = mixpanel.NewConnector(apiURL, token)
	}

	return c.mixpanel
}
