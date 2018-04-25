package metrics

import (
	"sync"

	"fmt"
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app/internal/metrics/data"
	"gitlab.com/swarmfund/psim/psim/app/internal/metrics/middlewares"
)

func (m *MetricsProvider) ServicesStatus(w http.ResponseWriter, r *http.Request) {
	if !m.IsDone() {
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	result := map[string]map[string]map[string]interface{}{}

	statusCode := http.StatusOK
	for name, val := range m.metrics {
		if !val.State() {
			statusCode = http.StatusServiceUnavailable
		}
		result[name] = val.GetRegister().GetAll()
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(result)
}

func (m *MetricsProvider) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(m.log),
		middlewares.ContentType("application/json"),
	)

	r.Get("/services", m.ServicesStatus)

	return r
}

type MetricsProvider struct {
	done    bool
	mutex   *sync.Mutex
	log     *logan.Entry
	config  Config
	metrics map[string]*data.Metrics
}

func New(log *logan.Entry, config Config) *MetricsProvider {
	return &MetricsProvider{
		mutex:   &sync.Mutex{},
		log:     log,
		config:  config,
		metrics: map[string]*data.Metrics{},
	}
}

func (m *MetricsProvider) Run() {
	r := m.Router()

	addr := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)
	if err := http.ListenAndServe(addr, r); err != nil {
		m.log.WithError(err).Error("Failed start app")
		return
	}
}

func (m *MetricsProvider) AddService(name string) *data.Metrics {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.metrics[name] = data.NewMetric()
	return m.metrics[name]
}

func (m *MetricsProvider) Done() {
	m.done = true
}

func (m MetricsProvider) IsDone() bool {
	return m.done
}

//CheckServicesHealth sends request to Keeper and check services status
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

	entry.WithFields(logan.F{"metrics": "check services"}).Info(response.Status)

	return nil
}
