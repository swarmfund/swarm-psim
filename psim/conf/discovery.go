package conf

import (
	discovery "gitlab.com/distributed_lab/discovery-go"
)

var (
	discoveryClient *discovery.Client
)

func (c *ViperConfig) Discovery() (*discovery.Client, error) {
	if discoveryClient == nil {
		var err error
		config := &discovery.ClientConfig{}
		v := c.viper.Sub("discovery")
		if v != nil {
			config.Env = v.GetString("env")
			config.Host = v.GetString("host")
			config.Port = v.GetInt("port")
		}
		discoveryClient, err = discovery.NewClient(config)
		if err != nil {
			return nil, err
		}
	}
	return discoveryClient, nil
}
