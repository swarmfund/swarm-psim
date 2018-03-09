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
	StatusCode  int
	apiResponse *apiResponse
}

func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// DEPRECATED - just doesn't work. Look directly into StatusCode instead
func (r *Response) IsPermanent() bool {
	if r.StatusCode == http.StatusTooManyRequests && r.apiResponse != nil {
		for _, e := range r.apiResponse.Errors {
			if e.IsPermanent {
				return true
			}
		}
		return false
	}
	return r.StatusCode >= 400 && r.StatusCode < 500
}

func (r *Response) Authenticated() bool {
	return r.StatusCode != http.StatusUnauthorized && r.StatusCode < 500
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

func (r Response) GetLoganFields() map[string]interface{} {
	return map[string]interface{} {
		"status_code": r.StatusCode,
		"body":        r.apiResponse,
	}
}
