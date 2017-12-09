package utils

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	discovery "gitlab.com/distributed_lab/discovery-go"
)

type AskResult int

const (
	AskResultSuccess AskResult = iota
	AskResultFailure
	AskResultPermanentFailure
)

type AskMeta struct {
	NeighborsAsked int
	Err            error
}

// following procedure assumes neighbors are honest and follow same protocol,
// both for request/response schema and interaction expectations
func AskNeighbors(service *discovery.Service, payload interface{}) (AskResult, *AskMeta) {
	meta := &AskMeta{}
	body, err := json.Marshal(&payload)
	if err != nil {
		meta.Err = errors.Wrap(err, "failed to marshal payload")
		return AskResultPermanentFailure, meta
	}
	neighbors, err := service.DiscoverNeighbors()
	if err != nil {
		meta.Err = errors.Wrap(err, "failed to discover neighbors")
		return AskResultPermanentFailure, meta
	}
	for _, neighbor := range neighbors {
		response, err := http.Post(
			neighbor.Address, "application/json", bytes.NewReader(body),
		)
		if err != nil {
			// silencing err here, treating connection error just if neighbor does not exists
			continue
		}
		defer response.Body.Close()
		switch response.StatusCode {
		case http.StatusOK:
			// neighbor agreed with request and it went ok
			return AskResultSuccess, meta
		case http.StatusBadRequest, http.StatusUnsupportedMediaType:
			// according to neighbor request was not valid,
			// assuming they are all running same protocol it's permanent issue
			// TODO add message to meta
			return AskResultPermanentFailure, meta
		case http.StatusExpectationFailed:
			// neighbor does not agree with request or just unable to verify it ATM
		default:
			// consider it neighbor error
		}
	}
	return AskResultFailure, meta
}
