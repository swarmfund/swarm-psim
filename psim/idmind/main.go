package idmind

import (
	"context"

	"reflect"

	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceIdentityMind, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx)

	var config Config
	err := figure.
		Out(&config).
		From(app.Config(ctx).GetRequired(conf.ServiceIdentityMind)).
		With(figure.BaseHooks, utils.ETHHooks, hooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out", logan.F{
			"service": conf.ServiceIdentityMind,
		})
	}

	horizonConnector := globalConfig.Horizon().WithSigner(config.Signer)

	horizonInfo, err := horizonConnector.Info()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Horizon info")
	}

	builder := xdrbuild.NewBuilder(horizonInfo.Passphrase, horizonInfo.TXExpirationPeriod)

	return NewService(
		log,
		config,
		horizonConnector.Listener(),
		horizonConnector.Blobs(),
		horizonConnector.Users(),
		horizonConnector.Documents(),
		newConnector(config.Connector),
		builder,
	), nil
}

var hooks = figure.Hooks{
	"idmind.ConnectorConfig": func(raw interface{}) (reflect.Value, error) {
		rawConnectorConfig, err := cast.ToStringMapE(raw)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "failed to cast provider to map[string]interface{}")
		}

		var config ConnectorConfig
		err = figure.
			Out(&config).
			From(rawConnectorConfig).
			With(figure.BaseHooks).
			Please()
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "failed to figure out EmailsConfig")
		}

		return reflect.ValueOf(config), nil
	},
}
