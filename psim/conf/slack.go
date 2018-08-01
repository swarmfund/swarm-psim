package conf

import (
	"net/url"

	"github.com/multiplay/go-slack"
	"github.com/multiplay/go-slack/webhook"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (c *ViperConfig) Slack() slack.Client {
	c.Lock()
	defer c.Unlock()

	if c.slackClient != nil {
		return c.slackClient
	}

	var config struct {
		WebhookURL *url.URL `fig:"webhook_url"`
	}

	err := figure.Out(&config).From(c.GetRequired("slack")).Please()
	if err != nil {
		panic(errors.Wrap(err, "Failed to parse bitcoin config entry"))
	}

	c.slackClient = webhook.New(config.WebhookURL.String())

	return c.slackClient
}
