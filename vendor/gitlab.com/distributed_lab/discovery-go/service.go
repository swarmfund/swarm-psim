package discovery

import (
	"time"

	consul "github.com/hashicorp/consul/api"
)

const (
	DefaultTTL             = 10 * time.Second
	DefaultDeregisterAfter = 1 * time.Minute
)

type Service struct {
	client       *Client
	registration *ServiceRegistration
}

func NewService(client *Client, registration *ServiceRegistration) *Service {
	return &Service{
		client:       client,
		registration: registration,
	}
}

func (s *Service) DiscoverNeighbors() ([]ServiceEntry, error) {
	result := []ServiceEntry{}
	services, _, err := s.client.consul.Health().Service(
		s.registration.Name, "", true, nil)
	if err != nil {
		return nil, err
	}
	for _, service := range services {
		if s.registration.ID == service.Service.ID {
			continue
		}
		serviceEntry := ServiceEntry{
			ID:      service.Service.ID,
			Address: service.Service.Address,
		}
		result = append(result, serviceEntry)
	}
	return result, nil
}

type ServiceRegistration struct {
	Name            string
	ID              string
	TTL             time.Duration
	DeregisterAfter time.Duration

	// DEPRECATED use addr
	Host string
	// DEPRECATED use addr
	Port int
	Addr string
}

func (r *ServiceRegistration) ConsulRegistration() *consul.AgentServiceRegistration {
	// TODO generate random ID if empty
	ttl := DefaultTTL
	deregisterAfter := DefaultDeregisterAfter
	if r.TTL.Nanoseconds() > 0 {
		ttl = r.TTL
	}
	if r.DeregisterAfter.Nanoseconds() > 0 {
		deregisterAfter = r.DeregisterAfter
	}
	return &consul.AgentServiceRegistration{
		Name:    r.Name,
		ID:      r.ID,
		Address: r.Host,
		Port:    r.Port,
		Check: &consul.AgentServiceCheck{
			TTL: ttl.String(),
			DeregisterCriticalServiceAfter: deregisterAfter.String(),
		},
	}
}
