package conf

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type MandrillConf struct {
	API       string
	Key       string
	FromEmail string
	FromName  string
}

const mandrillConfigKey = "mandrill"

func (c *ViperConfig) Mandrill() MandrillConf {
	c.Lock()
	defer c.Unlock()

	if c.mandrill != nil {
		return *c.mandrill
	}

	mandrill := new(MandrillConf)
	config := c.GetStringMap(mandrillConfigKey)

	if err := figure.Out(mandrill).From(config).Please(); err != nil {
		panic(errors.Wrap(err, "failed to figure out mandrill "))
	}

	c.mandrill = mandrill
	return *c.mandrill
}
