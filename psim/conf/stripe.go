package conf

import (
	"errors"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

var (
	stripeClient *client.API
)

func (c *ViperConfig) Stripe() (*client.API, error) {
	if stripeClient == nil {
		v := c.viper.Sub("stripe")
		secret := v.GetString("secret_key")
		if secret == "" {
			return nil, errors.New("secret_key is required")
		}

		stripeClient = &client.API{}
		stripeClient.Init(secret, nil)
		stripe.LogLevel = 0
	}
	return stripeClient, nil
}
