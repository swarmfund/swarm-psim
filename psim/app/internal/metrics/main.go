package metrics

import (
	"sync"

	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app/internal/data"
	"gitlab.com/swarmfund/psim/psim/app/internal/metrics/middlewares"
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

func New(log *logan.Entry, config Config) *Metrics {
	return &Metrics{
		mutex:            &sync.Mutex{},
		log:              log,
		config:           config,
		trackingServices: []data.Metered{},
	}
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

//CheckServicesHealth sends request to metrics and check services status
//it takes info about metrics from config file
func CheckHealth(configData map[string]interface{}, entry *logan.Entry) error {
	metricsInfo, err := NewConfig(configData)
	if err != nil {
		return errors.Wrap(err, "Failed to get metrics config info")
	}

	url := fmt.Sprintf("http://%s:%d/services", metricsInfo.Host, metricsInfo.Port)
	response, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "Failed to get response from metrics")
	}

	//check status code
	switch response.StatusCode {
	case http.StatusOK:
		entry.WithFields(logan.F{"metrics": "StatusOK"}).Info("all services were successfully launched ")
	case http.StatusPreconditionFailed:
		entry.WithFields(logan.F{"metrics": "StatusPreconditionFailed"}).Error("services not initialized yet")
	case http.StatusServiceUnavailable:
		entry.WithFields(logan.F{"metrics": "StatusServiceUnavailable"}).Error("some of the services could not start")
	default:
		entry.WithFields(logan.F{"metrics": "unrecognized status code"}).Error(fmt.Sprintf("received code %d", response.StatusCode))
	}

	return nil
}
