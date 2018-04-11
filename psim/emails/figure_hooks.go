package emails

import (
	"reflect"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var FigureHooks = figure.Hooks{
	"email.Config": func(raw interface{}) (reflect.Value, error) {
		rawEmails, err := cast.ToStringMapE(raw)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "failed to cast provider to map[string]interface{}")
		}

		var emails Config
		err = figure.
			Out(&emails).
			From(rawEmails).
			With(figure.BaseHooks).
			Please()
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "failed to figure out EmailsConfig")
		}

		return reflect.ValueOf(emails), nil
	},
}
