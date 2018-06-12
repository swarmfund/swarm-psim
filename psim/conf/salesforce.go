package conf

import (
	"net/url"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/salesforce"
)

// TODO use figure out

func (c *ViperConfig) Salesforce() *salesforce.Connector {
	c.Lock()
	defer c.Unlock()

	if c.salesforce == nil {
		v := c.viper.Sub("salesforce")
		if v == nil {
			panic("salesforce config entry is missing")
		}

		apiRawURL := v.GetString("api_url")
		apiURL, err := url.Parse(apiRawURL)
		if err != nil {
			panic(errors.Wrap(err, "failed to parse api url", logan.F{
				"api_url": apiRawURL,
			}))
		}
		secret := v.GetString("client_secret")
		id := v.GetString("client_id")

		c.salesforce, err = salesforce.NewConnector(apiURL, secret, id)
		if err != nil {
			panic(errors.Wrap(err, "failed to create connector"))
		}
	}

	return c.salesforce
}
