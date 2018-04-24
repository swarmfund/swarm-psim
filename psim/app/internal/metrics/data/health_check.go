package data

type HealthCheck struct {
	err error
}

func NewHealthCheck() *HealthCheck {
	return &HealthCheck{nil}
}

//no-op
func (h *HealthCheck) Check() {}

// Error returns the healthCheck's status, which will be nil if it is healthy.
func (h *HealthCheck) Error() error {
	return h.err
}

// Healthy marks the healthCheck as healthy.
func (h *HealthCheck) Healthy() {
	h.err = nil
}

// Unhealthy marks the healthCheck as unhealthy.  The error is stored and
// may be retrieved by the Error method.
func (h *HealthCheck) Unhealthy(err error) {
	h.err = err
}
