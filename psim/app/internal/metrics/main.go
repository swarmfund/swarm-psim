package metrics

import (
	"sync"

	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app/internal/data"
	"gitlab.com/swarmfund/psim/psim/app/internal/metrics/handlers"
	"gitlab.com/swarmfund/psim/psim/app/internal/metrics/middlewares"
	"gitlab.com/swarmfund/psim/psim/conf"
)

func (a *Metrics) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(a.log),
		middlewares.ContentType("application/json"),
		middlewares.Ctx(
			handlers.CtxServices(a.trackingServices),
			handlers.CtxMutex(a.mutex),
		),
	)

	r.Get("/services", handlers.ServicesStatus)

	return r
}

type Metrics struct {
	mutex            *sync.Mutex
	log              *logan.Entry
	config           Config
	trackingServices []data.MetricsKeeper
}

func New(log *logan.Entry, configData map[string]interface{}) (*Metrics, error) {
	config := Config{}

	err := figure.
		Out(&config).
		From(configData).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out", logan.F{
			"service": conf.ServiceMetrics,
		})
	}

	return &Metrics{
		mutex:            &sync.Mutex{},
		log:              log,
		config:           config,
		trackingServices: []data.MetricsKeeper{},
	}, nil
}

func (m *Metrics) Run() {
	r := m.Router()

	addr := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)
	if err := http.ListenAndServe(addr, r); err != nil {
		m.log.WithError(err).Error("Failed start app")
		return
	}
}

func (m *Metrics) Unlock() {
	m.mutex.Unlock()
}

func (m *Metrics) Lock() {
	m.mutex.Lock()
}

func (m *Metrics) AddService(service data.MetricsKeeper) {
	m.trackingServices = append(m.trackingServices, service)
}
