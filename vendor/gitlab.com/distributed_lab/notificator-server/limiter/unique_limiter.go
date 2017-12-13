package limiter

import (
	"gitlab.com/distributed_lab/notificator-server/q"
	"gitlab.com/distributed_lab/notificator-server/types"
)

type UniqueLimiter struct {
}

func NewUniqueLimiter() *UniqueLimiter {
	return &UniqueLimiter{}
}

func (l UniqueLimiter) Check(request *types.APIRequest) (*CheckResult, error) {
	hash := request.GetHash()
	result, err := q.Request().ByHash(hash)
	if err != nil {
		return nil, err
	}

	return &CheckResult{
		Success: result == nil,
		IsPermanent: result != nil,
	}, nil
}

