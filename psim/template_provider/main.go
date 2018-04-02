package template_provider

import (
	"context"

	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"

	"fmt"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func init() {
	app.RegisterService(conf.ServiceTemplateProvider, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	globalConfig := app.Config(ctx)

	//default
	api := TemplateAPI{
		Host: "localhost",
		Port: 2323,
	}
	err := figure.Out(&api).From(globalConfig.Get(templateAPI)).Please()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to figure out %s", templateAPI))
	}

	return New(app.Config(ctx).TemplateProvider(), app.Log(ctx), api, app.Config(ctx).Horizon()), nil
}
