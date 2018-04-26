package investready

import (
	"reflect"

	"time"

	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	Connector       ConnectorConfig `fig:"connector,required,required"`
	RedirectsConfig RedirectsConfig `fig:"redirects,required"`

	Source keypair.Address `fig:"source,required"`
	Signer keypair.Full    `fig:"signer,required" mapstructure:"signer"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"connector": c.Connector,
	}
}

type ConnectorConfig struct {
	URL          string        `fig:"url,required"`
	ClientID     string        `fig:"client_id,required"`
	ClientSecret string        `fig:"client_secret,required"`
	Timeout      time.Duration `fig:"timeout"`
}

func (c ConnectorConfig) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"url":     c.URL,
		"timeout": c.Timeout,
	}
}

type RedirectsConfig struct {
	Host    string        `fig:"host,required"`
	Port    int           `fig:"port,required"`
	Timeout time.Duration `fig:"timeout"`
}

func (c RedirectsConfig) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"host":    c.Host,
		"port":    c.Port,
		"timeout": c.Timeout,
	}
}

var hooks = figure.Hooks{
	"investready.ConnectorConfig": func(raw interface{}) (reflect.Value, error) {
		rawConnectorConfig, err := cast.ToStringMapE(raw)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "Failed to cast provider to map[string]interface{}")
		}

		var config ConnectorConfig
		err = figure.
			Out(&config).
			From(rawConnectorConfig).
			With(figure.BaseHooks).
			Please()
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "Failed to figure out ConnectorConfig")
		}

		return reflect.ValueOf(config), nil
	},
	"investready.RedirectsConfig": func(raw interface{}) (reflect.Value, error) {
		rawRedirectsConfig, err := cast.ToStringMapE(raw)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "Failed to cast provider to map[string]interface{}")
		}

		var config RedirectsConfig
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
