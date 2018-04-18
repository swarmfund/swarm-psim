package data

type Metered interface {
	//HealthCheck used for check health of different services,
	HealthCheck() bool
}
