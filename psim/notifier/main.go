package notifier

import (
	"context"
	"fmt"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/figure"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceOperationNotifier, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	cfg := &Config{}
	err := figure.Out(cfg).
		From(globalConfig.Get(conf.ServiceOperationNotifier)).
		With(figure.BaseHooks, utils.CommonHooks, CommonHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err,
			fmt.Sprintf("failed to figure out %s",
				conf.ServiceOperationNotifier))
	}

	sender, err := globalConfig.Notificator()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get horizon connector")
	}

	horizonConn, err := globalConfig.Horizon()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get horizon connector")
	}

	logger := app.Log(ctx)
	logger = logger.WithField("service", "notifier")
	logger.Info("Starting")

	return &Service{
		Config:  cfg,
		horizon: horizonConn,
		logger:  logger,
		sender:  sender,
	}, nil
}
