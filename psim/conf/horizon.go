package conf

import (
	"github.com/pkg/errors"
	"gitlab.com/tokend/horizon-connector"
)

var (
	horizonClient *horizon.Connector
)

func (c *ViperConfig) Horizon() (*horizon.Connector, error) {
	if horizonClient != nil {
		return horizonClient, nil
	}

	v := c.viper.Sub("horizon")
	var endpoint string
	if v != nil {
		endpoint = v.GetString("addr")
	}

	client, err := horizon.NewConnector(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init horizon connector")
	}

	horizonClient = client
	return horizonClient, nil
}
