package handlers

import (
	"net/http"
)

func ServicesStatus(w http.ResponseWriter, r *http.Request) {
	Mutex(r).Lock()
	defer Mutex(r).Unlock()

	services := Services(r)

	var healthyAmount int
	for _, service := range services {
		if service.HealthCheck() {
			healthyAmount++
		}
	}

	if (len(services) - healthyAmount) != 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
}
