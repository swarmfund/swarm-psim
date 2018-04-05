package template_provider

import (
	"context"

	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func init() {
	app.RegisterService(conf.ServiceTemplateProvider, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)

	//default
	api := Config{
		Host: "localhost",
		Port: 2323,
	}
	err := figure.Out(&api).From(globalConfig.Get(conf.ServiceTemplateProvider)).Please()
	if err != nil {
		return nil, errors.Wrap(err, "failed to figure out")
	}

	info, err := app.Config(ctx).Horizon().Info()
	if err != nil {
		app.Log(ctx).WithError(err).Error("failed to get horizon")
		return nil, errors.Wrap(err, "Failed to get horizon info")
	}

	return New(app.Config(ctx).S3(), app.Log(ctx), api, info, app.Config(ctx).Horizon()), nil
}
