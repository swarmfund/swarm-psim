package conf

import (
	"errors"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

func (c *ViperConfig) Stripe() (*client.API, error) {
	c.Lock()
	defer c.Unlock()

	if c.stripeClient == nil {
		v := c.viper.Sub("stripe")
		secret := v.GetString("secret_key")
		if secret == "" {
			return nil, errors.New("secret_key is required")
		}

		c.stripeClient = &client.API{}
		c.stripeClient.Init(secret, nil)
		stripe.LogLevel = 0
	}
	return c.stripeClient, nil
}
