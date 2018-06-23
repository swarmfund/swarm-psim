package listener

import (
	"reflect"
	"time"

	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var ConfigFigureHooks = figure.Hooks{
	"listener.Config": func(raw interface{}) (reflect.Value, error) {
		rawRedirectsConfig, err := cast.ToStringMapE(raw)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "Failed to cast provider to map[string]interface{}")
		}

		var config Config
		err = figure.
			Out(&config).
			From(rawRedirectsConfig).
			With(figure.BaseHooks).
			Please()
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "Failed to figure out RedirectsConfig")
		}

		return reflect.ValueOf(config), nil
	},
}

type Config struct {
	Host           string        `fig:"host,required"`
	Port           int           `fig:"port,required"`
	Timeout        time.Duration `fig:"timeout"`
	CheckSignature bool          `fig:"check_signature,required"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"host":            c.Host,
		"port":            c.Port,
		"timeout":         c.Timeout,
		"check_signature": c.CheckSignature,
	}
}
