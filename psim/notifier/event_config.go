package notifier

import (
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"reflect"
)

type EventConfig struct {
	Cursor uint64       `fig:"cursor"`
	Emails EmailsConfig `fig:"emails,required"`
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
