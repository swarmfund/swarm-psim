package supervisor

import (
	"fmt"

	"reflect"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/keypair"
	"gitlab.com/swarmfund/psim/figure"
	"gitlab.com/swarmfund/psim/psim/utils"
)

type Config struct {
	Host  string
	Port  int
	Pprof bool `fig:"pprof"`

	LeadershipKey string

	SignerKP   keypair.KP `fig:"signer" mapstructure:"signer"`
	ExchangeKP keypair.KP `fig:"exchange"`
}

func NewConfig(serviceName string) Config {
	return Config{
		LeadershipKey: fmt.Sprintf("service/%s/leader", serviceName),
		Host:          "localhost",
	}
}

var ConfigFigureHooks = figure.Merge(figure.BaseHooks, utils.CommonHooks,
	figure.Hooks{
		"supervisor.Config": func(raw interface{}) (reflect.Value, error) {
			result := Config{}
			err := figure.Out(&result).
				From(raw.(map[string]interface{})).
				With(figure.BaseHooks, utils.CommonHooks).
				Please()
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to figure out supervisor common")
			}
			return reflect.ValueOf(result), nil
		},
	})
