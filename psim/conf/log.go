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

const (
	NLinesAroundErrorPoint = 2

	defaultLogLevel = "warn"
)

var (
	defaultLog *logan.Entry
)

func (c *ViperConfig) Log() (*logan.Entry, error) {
	if defaultLog != nil {
		return defaultLog, nil
	}

	// TODO Consider creating LogConfig struct and adding parsing of the 'log' config block into this struct using viper.Unmarshal()
	logViper := c.viper.Sub("log")
	if logViper == nil {
		return nil, errors.New("Log config is required.")
	}

	entry := logan.New()
	level := logViper.GetString("level")
	if level == "" {
		level = defaultLogLevel
	}
	lvl, err := logan.ParseLevel(level)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse log level")
	}
	entry = entry.Level(lvl)

	// horizon submitter hook
	// FIXME Make horizon-connector use logrus Hooks, not logan and then ei will work and can be uncommented.
	//entry.AddLogrusHook(&horizon.TXFailedHook{})

	err = addSlackHook(logViper, entry)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to add Slack hook")
	}

	err = addSentryHook(logViper, entry)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to add Sentry hook")
	}

	// set log formatter
	formatter := logViper.GetString("formatter")
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
		return errors.Wrap(err, "Failed to create Sentry hook")
	}

	env := v.GetString("env")
	if env == "" {
		env = "unknown"
	}
	hook.SetEnvironment(env)

	hook.StacktraceConfiguration.Enable = true
	// TODO Consider using log level from config
	hook.StacktraceConfiguration.Level = logrus.ErrorLevel
	hook.StacktraceConfiguration.Context = NLinesAroundErrorPoint

	wrapperHook := sentryWrapperHook{
		SentryHook: hook,
	}

	entry.AddLogrusHook(&wrapperHook)
	return nil
}

type sentryWrapperHook struct {
	*logrus_sentry.SentryHook
}

func (h *sentryWrapperHook) Fire(entry *logrus.Entry) error {
	err, ok := entry.Data[logan.ErrorKey]
	if ok {
		entry.Data["raw_error"] = err
	}

	return h.SentryHook.Fire(entry)
}
