// request_monitor wakes up every SleepPeriod (see config),
// figures out which requests haven't been resolved in a specified period
// and how many requests of each type there are
// and submits this info to console via logger
package request_monitor

import (
	"context"

	"time"

	"fmt"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/notifications"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceRequestMonitor, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	config := Config{
		SleepPeriod: 1 * time.Minute,
	}
	err := figure.
		Out(&config).
		From(app.Config(ctx).GetRequired(conf.ServiceRequestMonitor)).
		With(figure.BaseHooks, RequestsHook, utils.ETHHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "failed to figure out")
	}

	var notifier Notifier
	notifier = &LoganNotifier{app.Log(ctx)}
	if config.EnableSlack {
		notifier = &SlackNotifier{app.Config(ctx).NotificationSender()}
	}

	return New(config, notifier, app.Log(ctx), app.Config(ctx).Horizon()), nil
}

type Notifier interface {
	Notify(fields map[string]interface{}, msg string) error
}

type SlackNotifier struct {
	slackSender *notifications.SlackSender
}

func (n *SlackNotifier) Notify(fields map[string]interface{}, msg string) error {
	for key, value := range fields {
		msg += fmt.Sprintf(", %s : %v", key, value)
	}

	err := n.slackSender.Send(msg)
	if err != nil {
		return errors.Wrap(err, "failed to send slack message")
	}

	return nil
}

type LoganNotifier struct {
	log *logan.Entry
}

func (n *LoganNotifier) Notify(fields map[string]interface{}, msg string) error {
	n.log.WithFields(fields).Error(msg)
	return nil
}
