package conf

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// Conf is a configuration for Mailgun Client.
type MailgunConf struct {
	Key       string
	PublicKey string
	Domain    string
	From      string
}

const mailgunConfigKey = "mailgun"

func (c *ViperConfig) Mailgun() MailgunConf {
	c.Lock()
	defer c.Unlock()

	if c.mailgun != nil {
		return *c.mailgun
	}

	mailgun := new(MailgunConf)
	config := c.GetStringMap(mailgunConfigKey)

	if err := figure.Out(mailgun).From(config).Please(); err != nil {
		panic(errors.Wrap(err, "failed to figure out mailgun"))
	}

	c.mailgun = mailgun
	return *c.mailgun
}
