package limiter

import (
	"gitlab.com/distributed_lab/notificator-server/q"
	"gitlab.com/distributed_lab/notificator-server/types"
)

//DirectlyUniqueLimiter is intended for requests that assume their token is unique.
type DirectlyUniqueLimiter struct {
}

func NewDirectlyUniqueLimiter() *DirectlyUniqueLimiter {
	return &DirectlyUniqueLimiter{}
}

func (l DirectlyUniqueLimiter) Check(request *types.APIRequest) (*CheckResult, error) {
	token := request.Token
	result, err := q.Request().ByToken(token)
	if err != nil {
		return nil, err
	}

	return &CheckResult{
		Success:     result == nil,
		IsPermanent: result != nil,
	}, nil
}
