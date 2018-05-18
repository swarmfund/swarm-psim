package data

import (
	"github.com/rcrowley/go-metrics"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Metrics struct {
	register    metrics.Registry
	healthCheck metrics.Healthcheck
}

func NewMetric() Metrics {
	metric := Metrics{
		register:    metrics.NewRegistry(),
		healthCheck: NewHealthCheck(),
	}

	//Set register health checker and set default value to unhealthy
	metric.register.GetOrRegister("health", metric.healthCheck)
	metric.Unhealthy(errors.New("services not initialize yet"))

	return metric
}

func (m *Metrics) Register(name string, value interface{}) {
	m.register.Register(name, value)
}

func (m *Metrics) Healthy() {
	m.healthCheck.Healthy()
}

func (m *Metrics) Unhealthy(err error) {
	m.healthCheck.Unhealthy(err)
}

func (m *Metrics) State() bool {
	return m.healthCheck.Error() == nil
}

func (m Metrics) GetRegister() metrics.Registry {
	return m.register
}
