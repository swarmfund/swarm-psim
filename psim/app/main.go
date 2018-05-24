package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/conf"
)

const (
	// TODO Custom types
	// DEPRECATED use Config getter instead
	CtxConfig = "config"
	// DEPRECATED use Log getter instead
	CtxLog                 = "log"
	ctxLog                 = CtxLog
	ctxConfig              = CtxConfig
	forceKillPeriodSeconds = 30
)

// DEPRECATED
type DeprecatedService func(ctx context.Context) error
type ServiceSetUp func(ctx context.Context) (Service, error)

var (
	servicesMu           = sync.RWMutex{}
	registeredServices   = map[string]DeprecatedService{}
	registerServiceSetUp = map[string]ServiceSetUp{}
)

type Service interface {
	Run(ctx context.Context)
}

func Config(ctx context.Context) conf.Config {
	v := ctx.Value(ctxConfig)
	if v == nil {
		panic("config context value expected")
	}
	return v.(conf.Config)
}

func Log(ctx context.Context) *logan.Entry {
	v := ctx.Value(ctxLog)
	if v == nil {
		panic("log context value expected")
	}
	return v.(*logan.Entry)
}

func RegisterService(name string, setup ServiceSetUp) {
	servicesMu.Lock()
	defer servicesMu.Unlock()
	if setup == nil {
		panic("service set up fn is nil")
	}
	if _, duplicated := registerServiceSetUp[name]; duplicated {
		panic(fmt.Sprintf("service already registered %s", name))
	}
	registerServiceSetUp[name] = setup
}

// DEPRECATED use RegisterService
func Register(name string, service DeprecatedService) {
	servicesMu.Lock()
	defer servicesMu.Unlock()
	if service == nil {
		panic("register service is nil")
	}
	if _, dup := registeredServices[name]; dup {
		panic(fmt.Sprintf("service already registered %s", name))
	}
	registeredServices[name] = service
}

type App struct {
	log    *logan.Entry
	config conf.Config
	ctx    context.Context
	cancel context.CancelFunc
}

func New(config conf.Config) (*App, error) {
	entry, err := config.Log()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get logger")
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		config: config,
		log:    entry.WithField("service", "app"),
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (app *App) Run() {
	servicesMu.Lock()
	defer servicesMu.Unlock()

	app.log.WithField("services_count", len(app.config.Services())).Info("Starting services.")
	wg := sync.WaitGroup{}

	ctx := context.WithValue(app.ctx, ctxConfig, app.config)
	// <!-- DEPRECATED
	for name, service := range registeredServices {
		if !app.isServiceEnabled(name) {
			continue
		}
		wg.Add(1)
		ohaigo := service
		ohaigoagain := name
		go func() {
			// TODO defer panic and die
			defer wg.Done()
			entry := app.log.WithField("service", ohaigoagain)
			ctx := context.WithValue(ctx, ctxLog, entry)
			if err := ohaigo(ctx); err != nil {
				entry.WithError(err).Error("died")
			}
		}()
	}
	// -->

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		app.log.WithField("signal", sig).Info("Received signal.")
		app.cancel()
		done := make(chan struct{})

		// Close done after wg finishes.
		go func() {
			defer close(done)
			wg.Wait()
		}()

		select {
		case <-done:
			app.log.Info("Clean exit.")
			os.Exit(0)
		case <-time.NewTimer(forceKillPeriodSeconds * time.Second).C:
			// FIXME
			app.log.Warn("Failed to perform shutdown gracefully, some services couldn't stop - stopping now without waiting.")
			os.Exit(1)
		}
	}()

	for name, setup := range registerServiceSetUp {
		if !app.isServiceEnabled(name) {
			continue
		}

		wg.Add(1)
		ohaigo := setup
		entry := app.log.WithField("service", name)
		go func() {
			defer func() {
				if rec := recover(); rec != nil {
					err := errors.FromPanic(rec)
					entry.WithStack(errors.WithStack(err)).WithError(err).Error("Service panicked during initialization or executing.")
				}
				wg.Done()
			}()

			ctx := context.WithValue(ctx, ctxLog, entry)
			service, err := ohaigo(ctx)
			if err != nil {
				// TODO Consider panicking here instead of Error log and return.
				entry.WithError(err).Error("App failed to set up service.")
				return
			}

			// TODO Pass another ctx here - just for cancelling.
			service.Run(ctx)
			entry.Warn("died")
		}()
	}

	wg.Wait()
	time.Sleep(1)
	os.Exit(0)
}

// IsServiceEnabled returns true, if service is among services in config.
func (app *App) isServiceEnabled(name string) bool {
	for _, enabled := range app.config.Services() {
		if name == enabled {
			return true
		}
	}
	return false
}

func IsCanceled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
