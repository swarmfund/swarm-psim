package conf

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type HTTPConf struct {
	Host           string
	Port           int
	AllowUntrusted bool
}

const httpConfigKey = "http"

func (c *ViperConfig) HTTP() HTTPConf {
	c.Lock()
	defer c.Unlock()

	if c.http != nil {
		return *c.http
	}

	http := new(HTTPConf)
	config := c.GetStringMap(httpConfigKey)

	if err := figure.Out(http).From(config).Please(); err != nil {
		panic(errors.Wrap(err, "failed to figure out http"))
	}

	c.http = http

	return *c.http
}
