package conf

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/notifications"
)

func (c *ViperConfig) NotificationSender() *notifications.SlackSender {
	c.Lock()
	defer c.Unlock()

	if c.notificationSender != nil {
		return c.notificationSender
	}
	conf := notifications.SlackConfig{}

	err := figure.
		Out(&conf).
		From(c.Get("notifications_slack")).
		With(figure.BaseHooks).
		Please()

	if err != nil {
		panic(errors.Wrap(err, "Failed to figure out notifications_slack"))
	}

	sender := notifications.NewSlackSender(conf)

	c.notificationSender = sender

	return c.notificationSender
}
