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
	"gitlab.com/swarmfund/psim/psim/app/internal/metrics/middlewares"
	"gitlab.com/swarmfund/psim/psim/conf"
)

func (m *Metrics) ServicesStatus(w http.ResponseWriter, r *http.Request) {
	if !m.Done() {
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	services := m.GetServices()

	for _, service := range services {
		if !service.HealthCheck() {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (m *Metrics) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(m.log),
		middlewares.ContentType("application/json"),
	)

	r.Get("/services", m.ServicesStatus)

	return r
}

type Metrics struct {
	done             bool
	mutex            *sync.Mutex
	log              *logan.Entry
	config           Config
	trackingServices []data.Metered
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
		trackingServices: []data.Metered{},
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

func (m *Metrics) AddService(service data.Metered) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.trackingServices = append(m.trackingServices, service)
}
func (m *Metrics) GetServices() []data.Metered {
	return m.trackingServices
}

func (m *Metrics) SetDone(isDone bool) {
	m.done = isDone
}

func (m Metrics) Done() bool {
	return m.done
}
