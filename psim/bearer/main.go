// Package bearer is a runner which periodically
// submit some operation(s) to the Horizon.
package bearer

import (
	"context"

	"time"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.RegisterService(conf.ServiceBearer, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	config := Config{
		AbnormalPeriod: 1 * time.Minute,
		SleepPeriod:    5 * time.Minute,
	}
	err := figure.
		Out(&config).
		From(app.Config(ctx).Get(conf.ServiceBearer)).
		With(figure.BaseHooks, utils.ETHHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "failed to figure out")
	}

	horizon := app.Config(ctx).Horizon()
	info, err := horizon.Info()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get horizon info")
	}
	builder := xdrbuild.NewBuilder(info.Passphrase, info.TXExpirationPeriod)
	helper := NewCheckSalesStateHelper(app.Config(ctx).Horizon(), builder, config.Source, config.Signer)

	return New(config, app.Log(ctx), helper), nil
}
