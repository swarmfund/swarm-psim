package limiter

import (
	"time"

	"gitlab.com/distributed_lab/notificator-server/types"
)

type CheckResult struct {
	Success     bool
	IsPermanent bool
	RetryIn     *time.Duration
}

type Limiter interface {
	Check(request *types.APIRequest) (*CheckResult, error)
}

func ParseLimiter(raw map[string]interface{}) Limiter {
	switch raw["type"] {
	case "window":
		return NewWindowLimiter(raw)
	case "unique":
		return NewUniqueLimiter()
	case "directly-unique":
		return NewDirectlyUniqueLimiter()
	default:
		panic("unknown limiter type")
	}
}
