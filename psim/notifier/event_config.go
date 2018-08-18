package notifier

import (
	"reflect"

	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type EventConfig struct {
	Disabled bool         `fig:"disabled"`
	Cursor   uint64       `fig:"cursor"`
	Emails   EmailsConfig `fig:"emails,required"`
}

func (c EventConfig) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"cursor": c.Cursor,
		"emails": c.Emails,
	}
}

var EventHooks = figure.Hooks{
	"notifier.EventConfig": func(raw interface{}) (reflect.Value, error) {
		rawEvent, err := cast.ToStringMapE(raw)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "failed to cast provider to map[string]interface{}")
		}

		var disabledConfig struct {
			Disabled bool `fig:"disabled"`
		}

		err = figure.Out(&disabledConfig).From(rawEvent).Please()

		if disabledConfig.Disabled {
			return reflect.ValueOf(EventConfig{Disabled: true}), nil
		}

		var event EventConfig
		err = figure.
			Out(&event).
			From(rawEvent).
			With(figure.BaseHooks, EmailsHooks).
			Please()
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "failed to figure out EventConfig")
		}

		return reflect.ValueOf(event), nil
	},
}
