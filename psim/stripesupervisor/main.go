package stripesupervisor

import (
	"context"

	"fmt"

	"github.com/stripe/stripe-go/client"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/figure"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/supervisor"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	setupFn := func(ctx context.Context) (utils.Service, error) {
		globalConfig := app.Config(ctx)

		config := supervisor.NewConfig(conf.ServiceStripeSupervisor)
		err := figure.
			Out(&config).
			From(globalConfig.Get(conf.ServiceStripeSupervisor)).
			With(figure.BaseHooks, utils.CommonHooks).
			Please()
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to figure out %s", conf.ServiceStripeSupervisor))
		}

		commonSupervisor, err := supervisor.InitNew(ctx, conf.ServiceStripeSupervisor, config)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to init common supervisor")
		}

		stripeClient, err := globalConfig.Stripe()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get stripe client")
		}

		return newService(commonSupervisor, stripeClient), nil
	}

	app.RegisterService(conf.ServiceStripeSupervisor, setupFn)
}

// Service implements utils.Service interface, it supervises Stripe transactions
// and send CoinEmissionRequests to Horizon if arrived Charge detected.
//
// Service uses supervisor.Service for common for supervisors logic, such as Leadership and Profiling.
type Service struct {
	*supervisor.Service

	stripe *client.API
}

func newService(commonSupervisor *supervisor.Service, stripe *client.API) *Service {
	result := &Service{
		Service: commonSupervisor,

		stripe: stripe,
	}

	result.AddRunner(result.processStripeHistory)

	return result
}
