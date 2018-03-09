package conf

import (
	"fmt"
	"reflect"

	"github.com/Sirupsen/logrus"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const logConfigKey = "log"

var (
	logLevelHook = figure.Hooks{
		"logrus.Level": func(value interface{}) (reflect.Value, error) {
			switch v := value.(type) {
			case string:
				lvl, err := logrus.ParseLevel(v)
				if err != nil {
					return reflect.Value{}, errors.Wrap(err, "failed to parse log level")
				}
				return reflect.ValueOf(lvl), nil
			case nil:
				return reflect.ValueOf(nil), nil
			default:
				return reflect.Value{}, fmt.Errorf("unsupported conversion from %T", value)
			}
		},
	}
)

func (c *ViperConfig) Log() *logrus.Logger {
	if c.log != nil {
		return c.log
	}

	var config struct {
		Level logrus.Level
	}

	err := figure.
		Out(&config).
		With(figure.BaseHooks, logLevelHook).
		From(c.GetStringMap(logConfigKey)).
		Please()
	if err != nil {
		panic(errors.Wrap(err, "failed to figure out log"))
	}

	c.log = logrus.New()
	c.log.Level = logrus.Level(config.Level)

	return c.log
}
