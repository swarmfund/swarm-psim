package limiter

import (
	"time"
	"gitlab.com/distributed_lab/notificator-server/q"
	"gitlab.com/distributed_lab/notificator-server/types"
)

type WindowLimiter struct {
	Requests int
	Interval time.Duration
}

func NewWindowLimiter(raw map[string]interface{}) *WindowLimiter {
	interval, err := time.ParseDuration(raw["interval"].(string))
	if err != nil {
		panic(err)
	}
	return &WindowLimiter{
		Requests: raw["requests"].(int),
		Interval: interval,
	}
}

func (l WindowLimiter) Check(request *types.APIRequest) (*CheckResult, error) {
	tillNextWindow, err := q.Request().NextWindow(request.Type, request.Token, l.Requests, l.Interval)

	if err != nil {
		return nil, err
	}

	return &CheckResult{
		Success: tillNextWindow == nil,
		IsPermanent: false,
		RetryIn: tillNextWindow,
	}, nil
}
