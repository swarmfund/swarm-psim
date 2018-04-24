package data

import "github.com/rcrowley/go-metrics"

type Metrics struct {
	register    metrics.Registry
	healthCheck metrics.Healthcheck
}

func NewMetric() *Metrics {
	register := metrics.NewRegistry()
	healthCheck := NewHealthCheck()

	register.GetOrRegister("health", healthCheck)

	return &Metrics{
		register:    register,
		healthCheck: healthCheck,
	}
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

func (m Metrics) GetRegister() metrics.Registry {
	return m.register
}
