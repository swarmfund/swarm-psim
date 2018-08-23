package supervisor

import (
	"fmt"

	"reflect"

	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/externalsystems/derive"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	Host  string
	Port  int
	Pprof bool `fig:"pprof"`

	LeadershipKey string

	SignerKP   keypair.Full    `fig:"signer" mapstructure:"signer"`
	ExchangeKP keypair.Address `fig:"exchange"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"host": c.Host,
		"port": c.Port,

		"leadership_key": c.LeadershipKey,
	}
}

func NewConfig(serviceName string) Config {
	return Config{
		LeadershipKey: fmt.Sprintf("service/%s/leader", serviceName),
		Host:          "localhost",
	}
}

var DLFigureHooks = figure.Hooks{
	"derive.NetworkType": func(raw interface{}) (reflect.Value, error) {
		i, err := cast.ToInt32E(raw)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "int32 cast failed")
		}
		return reflect.ValueOf(derive.NetworkType(i)), nil
	},
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
}
