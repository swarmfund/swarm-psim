// Package bearer is a runner which periodically
// submit some operation(s) to the Horizon.
package bearer

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
	app.RegisterService(conf.ServiceBearer, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)
	log := app.Log(ctx).WithField("service", conf.ServiceBearer)

	config := Config{}

	err := figure.
		Out(&config).
		From(globalConfig.GetRequired(conf.ServiceBearer)).
		With(figure.BaseHooks, utils.CommonHooks).
		Please()

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to figure out %s", conf.ServiceBearer))
	}

	hConn := globalConfig.Horizon()
	checkSalesStateHelper := NewCheckSalesStateHelper(hConn, config)
	return New(config, log, checkSalesStateHelper), nil
}
