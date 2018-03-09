package conf

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type PostageConf struct {
	Key  string
	From string
}

const postageConfigKey = "postage"

func (c *ViperConfig) Postage() PostageConf {
	c.Lock()
	defer c.Unlock()

	if c.postage != nil {
		return *c.postage
	}

	postage := new(PostageConf)
	config := c.GetStringMap(postageConfigKey)

	if err := figure.Out(postage).From(config).Please(); err != nil {
		panic(errors.Wrap(err, "failed to figure out postage"))
	}
	c.postage = postage

	return *c.postage
}
