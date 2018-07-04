package request_monitor

import (
	"context"

	"time"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceRequestMonitor, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	config := Config{
		AbnormalPeriodMin: 1 * time.Minute,
		AbnormalPeriodMax: 10 * time.Minute,
		SleepPeriod: 1 * time.Minute,
	}
	err := figure.
		Out(&config).
		From(app.Config(ctx).Get(conf.ServiceRequestMonitor)).
		With(figure.BaseHooks, utils.ETHHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "failed to figure out")
	}

	return New(config, app.Log(ctx), app.Config(ctx).Horizon()), nil
}
