package limiter

import (
	"gitlab.com/distributed_lab/notificator-server/types"
	"time"
)

type CheckResult struct {
	Success     bool
	IsPermanent bool
	RetryIn     *time.Duration
}

type Limiter interface {
	Check(request *types.APIRequest) (*CheckResult, error)
}
