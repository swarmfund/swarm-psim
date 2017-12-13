package conf

import (
	"github.com/evalphobia/logrus_sentry"
	"github.com/multiplay/go-slack/chat"
	"github.com/multiplay/go-slack/lrhook"
	"github.com/sirupsen/logrus"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"github.com/spf13/viper"
)

const NLinesAroundErrorPoint = 2

var (
	defaultLog *logan.Entry
)

func (c *ViperConfig) Log() (*logan.Entry, error) {
	if defaultLog != nil {
		return defaultLog, nil
	}

	v := c.viper.Sub("log")
	if v == nil {
		return nil, errors.New("log config is required")
	}

	entry := logan.New()
	level := v.GetString("level")
	if level == "" {
		level = "warn"
	}
	lvl, err := logan.ParseLevel(level)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse level")
	}
	entry.Level(lvl)

	// horizon submitter hook
	// FIXME Make horizon-connector use logrus Hooks, not logan and then ei will work and can be uncommented.
	//entry.AddLogrusHook(&horizon.TXFailedHook{})

	err = addSlackHook(v, entry)
	if err != nil {
		return nil, err
	}

	err = addSentryHook(v, entry)
	if err != nil {
		return nil, err
	}

	// set log formatter
	formatter := v.GetString("formatter")
	if formatter != "" {
		switch formatter {
		case "json":
			entry.Formatter(logan.JSONFormatter)
		}
	}

	defaultLog = entry
	return defaultLog, nil
}

func addSlackHook(v *viper.Viper, entry *logan.Entry) error {
	webhook := v.GetString("slack_webhook")
	if webhook == "" {
		return nil
	}

	slackLevel := v.GetString("slack_level")
	if slackLevel == "" {
		slackLevel = "error"
	}
	slackLvl, err := logrus.ParseLevel(slackLevel)
	if err != nil {
		return errors.Wrap(err, "failed to parse slack level")
	}

	channel := v.GetString("slack_channel")
	if channel == "" {
		return errors.New("slack_channel is required")
	}

	cfg := lrhook.Config{
		MinLevel: slackLvl,
		Message: chat.Message{
			Channel:   channel,
			IconEmoji: ":glitch_crab:",
		},
	}

	h := lrhook.New(cfg, webhook)

	entry.AddLogrusHook(h)
	return nil
}

func addSentryHook(v *viper.Viper, entry *logan.Entry) error {
	sentry := v.GetString("sentry_dsn")
	if sentry == "" {
		return nil
	}

	hook, err := logrus_sentry.NewSentryHook(sentry, []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to initialize sentry")
	}

	hook.StacktraceConfiguration.Enable = true
	hook.StacktraceConfiguration.Level = logrus.ErrorLevel
	hook.StacktraceConfiguration.Context = NLinesAroundErrorPoint

	entry.AddLogrusHook(hook)
	return nil
}
