package conf

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type TwilioConf struct {
	SID        string
	Token      string
	FromNumber string
}

const twilioConfigKey = "twilio"

func (c *ViperConfig) Twilio() TwilioConf {
	c.Lock()
	defer c.Unlock()

	if c.twilio != nil {
		return *c.twilio
	}

	twilio := new(TwilioConf)
	config := c.GetStringMap(twilioConfigKey)

	if err := figure.Out(twilio).From(config).Please(); err != nil {
		panic(errors.Wrap(err, "failed to figure out twilio"))
	}

	c.twilio = twilio

	return *c.twilio
}
