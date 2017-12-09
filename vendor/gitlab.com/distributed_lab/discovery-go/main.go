package discovery

import (
	"fmt"

	"time"

	consul "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

const (
	DefaultEnv  = "default"
	DefaultHost = "localhost"
	DefaultPort = 8500
)

type Error struct {
	Fatal bool
	Err   error
}

func (e Error) Error() string {
	return e.Err.Error()
}

type Client struct {
	consul *consul.Client
	config *ClientConfig
}

func NewClient(config *ClientConfig) (*Client, error) {
	config = config.WithDefaults()
	consulClient, err := consul.NewClient(&consul.Config{
		Address: fmt.Sprintf("http://%s:%d", config.Host, config.Port),
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		consul: consulClient,
		config: config,
	}, nil
}

func (c *Client) Config() *ClientConfig {
	if c.config == nil {
		c.config = c.config.WithDefaults()
	}
	return c.config
}

func (c *Client) RegisterServiceSync(service *Service) error {
	err := c.consul.Agent().ServiceRegister(service.registration.ConsulRegistration())
	if err != nil {
		return errors.Wrap(err, "failed to register service")
	}
	ticker := time.NewTicker(5 * time.Second)
	serviceKey := fmt.Sprintf("service:%s", service.registration.ID)
	for ; true; <-ticker.C {
		err := c.consul.Agent().UpdateTTL(serviceKey, "OK", consul.HealthPassing)
		if err != nil {
			return errors.Wrap(err, "failed to update ttl")
		}
	}
	panic("run out of ticks")
}

func (c *Client) EnsureServiceRegistered(service *Service) (bool, error) {
	registration := service.registration.ConsulRegistration()
	// it's a redundant request every time, but who cares while it works
	err := c.consul.Agent().ServiceRegister(registration)
	if err != nil {
		return false, errors.Wrap(err, "failed to register")
	}
	err = c.consul.Agent().UpdateTTL(
		fmt.Sprintf("service:%s", service.registration.ID), "OK", consul.HealthPassing)
	if err != nil {
		return false, errors.Wrap(err, "failed to update ttl")
	}
	return true, nil
}

func (c *Client) DeregisterService(service *Service) (bool, error) {
	err := c.consul.Agent().ServiceDeregister(service.registration.ID)
	if err != nil {
		return false, errors.Wrap(err, "failed to deregister")
	}
	return true, nil
}

type ServiceEntry struct {
	ID      string
	Address string
}

func (c *Client) DiscoverService(service string) ([]ServiceEntry, error) {
	result := []ServiceEntry{}
	services, _, err := c.consul.Health().Service(service, "", true, nil)
	for _, service := range services {
		serviceEntry := ServiceEntry{
			ID:      service.Service.ID,
			Address: service.Service.Address,
		}
		result = append(result, serviceEntry)
	}
	return result, err
}

func (c *Client) Consul() *consul.Client {
	return c.consul
}

// kv should have session already
func (c *Client) TryAcquire(kv *KVPair) (bool, error) {
	ok, _, err := c.consul.KV().Acquire(kv.toConsul(), nil)
	return ok, err
}

func (c *Client) Acquire(kv KVPair) (*KVPair, error) {
	if kv.Session == nil {
		session, err := NewSession(c)
		if err != nil {
			return nil, err
		}
		kv.Session = session
		kv.Session.EndlessRenew()
	}

	for {
		ok, err := c.TryAcquire(&kv)
		if err != nil {
			return nil, err
		}
		if ok {
			return &kv, nil
		}
		time.Sleep(5 * time.Second)
	}
}

func (c *Client) Get(key string) (*KVPair, error) {
	kv, _, err := c.consul.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}
	if kv == nil {
		return nil, nil
	}
	return NewKV(c, kv), nil
}

func (c *Client) Set(kv KVPair) error {
	ckv := kv.toConsul()
	if _, err := c.consul.KV().Put(ckv, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetOrSet(kv KVPair) (*KVPair, error) {
	ckv, _, err := c.consul.KV().Get(kv.Key, nil)
	if err != nil {
		return nil, err
	}

	if ckv != nil {
		return NewKV(c, ckv), nil
	}

	ckv = kv.toConsul()
	ok, _, err := c.consul.KV().CAS(ckv, nil)
	if ok {
		return &kv, nil
	}

	ckv, _, err = c.consul.KV().Get(kv.Key, nil)
	if err != nil {
		return nil, err
	}

	return NewKV(c, ckv), nil
}

func (c *Client) Service(registration *ServiceRegistration) *Service {
	return &Service{
		client:       c,
		registration: registration,
	}
}
