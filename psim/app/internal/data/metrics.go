package data

type MetricsKeeper interface {
	//HealthCheck used for check health of different services,
	HealthCheck() bool
}
