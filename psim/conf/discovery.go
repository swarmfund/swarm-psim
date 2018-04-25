package conf

import (
	discovery "gitlab.com/distributed_lab/discovery-go"
)

func (c *ViperConfig) Discovery() *discovery.Client {
	c.Lock()
	defer c.Unlock()

	if c.discoveryClient == nil {
		var err error
		config := &discovery.ClientConfig{}
		v := c.viper.Sub("discovery")
		if v != nil {
			config.Env = v.GetString("env")
			config.Host = v.GetString("host")
			config.Port = v.GetInt("port")
		}
		c.discoveryClient, err = discovery.NewClient(config)
		if err != nil {
			panic(err)
		}
	}
	return c.discoveryClient
}
