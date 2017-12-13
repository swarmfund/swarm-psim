package notificator

import (
	"net/http"
	"time"
)

type apiResponse struct {
	Errors []apiError `json:"errors"`
}

type apiError struct {
	IsPermanent bool   `json:"is_permanent"`
	RetryIn     *int64 `json:"retry_in"`
}

type Response struct {
	statusCode  int
	apiResponse *apiResponse
}

func (r *Response) IsSuccess() bool {
	return r.statusCode >= 200 && r.statusCode < 300
}

func (r *Response) IsPermanent() bool {
	if r.statusCode == http.StatusTooManyRequests && r.apiResponse != nil {
		for _, e := range r.apiResponse.Errors {
			if e.IsPermanent {
				return true
			}
		}
		return false
	}
	return r.statusCode >= 400 && r.statusCode < 500
}

func (r *Response) Authenticated() bool {
	return r.statusCode != http.StatusUnauthorized && r.statusCode < 500
}

func (r *Response) RetryIn() *time.Duration {
	// TODO finding max would be nice
	if r.apiResponse != nil {
		for _, e := range r.apiResponse.Errors {
			if e.RetryIn != nil {
				retryIn := time.Second * time.Duration(*e.RetryIn)
				return &retryIn
			}
		}
	}
	return nil
}
