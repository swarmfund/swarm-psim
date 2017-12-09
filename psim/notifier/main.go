package notifier

import (
	"context"
	"fmt"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/psim/figure"
	"gitlab.com/tokend/psim/psim/app"
	"gitlab.com/tokend/psim/psim/conf"
	"gitlab.com/tokend/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceOperationNotifier, setupFn)
}

func setupFn(ctx context.Context) (utils.Service, error) {
	globalConfig := ctx.Value(app.CtxConfig).(conf.Config)
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

	logger := ctx.Value(app.CtxLog).(*logan.Entry)
	logger = logger.WithField("service", "notifier")
	logger.Info("Starting")

	ctxL, cancel := context.WithCancel(ctx)
	return &Service{
		Config:  cfg,
		horizon: horizonConn,
		logger:  logger,
		sender:  sender,
		ctx:     ctxL,
		cancel:  cancel,
	}, nil
}
