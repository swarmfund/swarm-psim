package conf

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/psim/psim/notifications"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var (
	notificationSender *notifications.SlackSender
)

func (c *ViperConfig) NotificationSender() *notifications.SlackSender {
	if notificationSender != nil {
		return notificationSender
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

	notificationSender = sender

	return sender
}
