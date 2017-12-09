package server

import (
	"gitlab.com/distributed_lab/notificator-server/limiter"
	"gitlab.com/distributed_lab/notificator-server/q"
	"gitlab.com/distributed_lab/notificator-server/types"
	"time"
)

type DispatchResultReasonType int

const (
	DispatchResultSuccess DispatchResultReasonType = iota
	DispatchResultUnknownType
	DispatchResultLimitExceeded
)

type RequestDispatcher struct {
}

type DispatchResult struct {
	Type        DispatchResultReasonType
	IsPermanent bool
	RetryIn     *time.Duration
}

func NewRequestDispatcher() *RequestDispatcher {
	return &RequestDispatcher{}
}

func (d *RequestDispatcher) Dispatch(apiRequest *types.APIRequest) (*DispatchResult, error) {
	requestsConf := GetRequestsConf()
	requestType, ok := requestsConf.Get(apiRequest.Type)

	if !ok {
		return &DispatchResult{Type: DispatchResultUnknownType}, nil
	}

	checkResults := []*limiter.CheckResult{}
	for _, limiter := range requestType.Limiters {
		checkResult, err := limiter.Check(apiRequest)
		if err != nil {
			return nil, err
		}

		checkResults = append(checkResults, checkResult)
	}

	dispatchResult := &DispatchResult{Type: DispatchResultSuccess}
	for _, checkResult := range checkResults {
		if checkResult.Success {
			continue
		}

		dispatchResult.Type = DispatchResultLimitExceeded

		if checkResult.IsPermanent {
			dispatchResult.IsPermanent = true
			break
		}

		if dispatchResult.RetryIn == nil || *dispatchResult.RetryIn < *checkResult.RetryIn {
			dispatchResult.RetryIn = checkResult.RetryIn
		}
	}

	if dispatchResult.Type == DispatchResultSuccess {
		err := q.Request().Insert(types.NewRequest(requestType.ID, requestType.Priority, apiRequest.PayloadString.Raw(), apiRequest.Token, apiRequest.GetHash()))
		if err != nil {
			return nil, err
		}
	}

	return dispatchResult, nil
}
