package investready

import (
	"reflect"

	"time"

	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/keypair"
	"gitlab.com/swarmfund/psim/psim/listener"
)

type Config struct {
	Connector       ConnectorConfig `fig:"connector,required,required"`
	RedirectsConfig listener.Config `fig:"redirects,required"`

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
	RedirectURI  string        `fig:"redirect_uri,required"`
	Timeout      time.Duration `fig:"timeout"` // optional
}

func (c ConnectorConfig) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"url":          c.URL,
		"client_id":    c.ClientID,
		"redirect_uri": c.RedirectURI,
		"timeout":      c.Timeout,
	}
}

var hooks = figure.Merge(figure.Hooks{
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
	},
	listener.ConfigFigureHooks,
)
