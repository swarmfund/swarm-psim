package bitcoin

import (
	"reflect"

	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// ConnectorConfig is structure to parse config for NodeConnector into.
type ConnectorConfig struct {
	Node           NodeConfig `fig:"node,required"`
	Testnet        bool       `fig:"testnet"`
	RequestTimeout int        `fig:"request_timeout_s"`
}

type NodeConfig struct {
	Host    string `fig:"host,required"`
	Port    int    `fig:"port,required"`
	AuthKey string `fig:"auth_key,required"`
}

func (c NodeConfig) GetLoganFields() map[string]interface{} {
	return map[string]interface{} {
		"host": c.Host,
		"post": c.Port,
	}
}

var FigureHooks = figure.Hooks{
	"bitcoin.NodeConfig": func(raw interface{}) (reflect.Value, error) {
		rawNode, err := cast.ToStringMapE(raw)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "Failed to cast NodeConfig to map[string]interface{}")
		}

		var nodeConfig NodeConfig
		err = figure.
			Out(&nodeConfig).
			From(rawNode).
			With(figure.BaseHooks).
			Please()
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "Failed to figure out NodeConfig")
		}

		return reflect.ValueOf(nodeConfig), nil
	},
}
