package template_provider

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/internal"
	"gitlab.com/tokend/go/doorman"
)

func init() {
	app.RegisterService(conf.ServiceTemplateProvider, setupFn)
}

func setupFn(ctx context.Context) (app.Service, error) {
	config, err := NewConfig(app.Config(ctx).Get(conf.ServiceTemplateProvider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to init config")
	}

	horizon := app.Config(ctx).Horizon()
	infoer := internal.NewLazyInfo(ctx, app.Log(ctx), horizon.System())
	doorman := doorman.New(config.SkipSignatureCheck, horizon.Accounts())

	return New(
		app.Config(ctx).S3(),
		app.Log(ctx),
		config,
		infoer,
		doorman,
	), nil
}
