package conf

import (
	"github.com/spf13/viper"
	"gitlab.com/distributed_lab/notificator-server/limiter"
	"gitlab.com/distributed_lab/notificator-server/types"
)

const (
	_ types.RequestTypeID = iota
	RequestTypeDummy
	RequestTypeUniqueDummy
	RequestTypeUniqueEmail
	RequestTypeUniqueSMS
)

type RequestType struct {
	ID       types.RequestTypeID
	Priority int
	Worker   string
	Limiters []limiter.Limiter
}

type RequestsConf struct {
	types map[types.RequestTypeID]RequestType
}

func (c *RequestsConf) Get(t types.RequestTypeID) (RequestType, bool) {
	value, ok := c.types[t]
	return value, ok
}

const requestsConfigKey = "requests"

func (c *ViperConfig) Requests() RequestsConf {
	c.Lock()
	defer c.Unlock()

	if c.requests != nil {
		return *c.requests
	}

	requestConf := map[string]struct {
		ID       types.RequestTypeID
		Priority int
		Limiters []map[string]interface{}
		Worker   string
	}{}

	err := viper.UnmarshalKey(requestsConfigKey, &requestConf)
	if err != nil {
		panic("failed to parse requests")
	}

	requestTypes := map[types.RequestTypeID]RequestType{}

	for _, v := range requestConf {

		requestType := RequestType{
			ID:       v.ID,
			Priority: v.Priority,
			Worker:   v.Worker,
			Limiters: make([]limiter.Limiter, len(v.Limiters)),
		}

		for i, l := range v.Limiters {
			requestLimiter := limiter.ParseLimiter(l)
			requestType.Limiters[i] = requestLimiter
		}
		requestTypes[v.ID] = requestType
	}

	c.requests = &RequestsConf{
		types: requestTypes,
	}

	return *c.requests
}
