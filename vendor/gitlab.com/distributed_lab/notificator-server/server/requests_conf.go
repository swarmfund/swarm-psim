package server

import (
	"github.com/spf13/viper"
	"gitlab.com/distributed_lab/notificator-server/limiter"
	"gitlab.com/distributed_lab/notificator-server/types"
)

type RequestsConf struct {
	types map[types.RequestTypeID]RequestType
}

func (c *RequestsConf) Get(t types.RequestTypeID) (RequestType, bool) {
	value, ok := c.types[t]
	return value, ok
}

func (c RequestsConf) Namespace() string {
	return "requests"
}

func (c RequestsConf) Init(v *viper.Viper) error {
	conf := map[string]struct {
		ID       types.RequestTypeID
		Priority int
		Limiters []map[string]interface{}
		Worker   string
	}{}

	err := viper.UnmarshalKey("requests", &conf)
	if err != nil {
		panic("failed to parse conf")
	}

	types := map[types.RequestTypeID]RequestType{}

	for _, v := range conf {
		requestType := RequestType{
			ID:       v.ID,
			Priority: v.Priority,
			Worker:   v.Worker,
			Limiters: make([]limiter.Limiter, len(v.Limiters)),
		}
		for i, l := range v.Limiters {
			limiter := parseLimiter(l)
			requestType.Limiters[i] = limiter
		}
		types[v.ID] = requestType
	}

	requestsConf = &RequestsConf{
		types: types,
	}
	return nil
}

var requestsConf *RequestsConf

func GetRequestsConf() *RequestsConf {
	return requestsConf
}

func parseLimiter(raw map[string]interface{}) limiter.Limiter {
	switch raw["type"] {
	case "window":
		return limiter.NewWindowLimiter(raw)
	case "unique":
		return limiter.NewUniqueLimiter()
	case "directly-unique":
		return limiter.NewDirectlyUniqueLimiter()
	default:
		panic("unknown limiter type")
	}
}
