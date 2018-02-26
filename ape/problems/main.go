package problems

import (
	"fmt"
	"net/http"

	"github.com/google/jsonapi"
)

// 400
func BadRequest(detail string) *jsonapi.ErrorObject {
	result := &jsonapi.ErrorObject{
		Title:  "Bad Request",
		Status: fmt.Sprintf("%d", http.StatusBadRequest),
	}

	if len(detail) > 0 {
		result.Detail = detail
	} else {
		result.Detail = "Your request was invalid in some way."
	}

	return result
}

// 403
func Forbidden(detail string) *jsonapi.ErrorObject {
	result := &jsonapi.ErrorObject{
		Title:  "Forbidden",
		Status: fmt.Sprintf("%d", http.StatusForbidden),
	}

	if len(detail) > 0 {
		result.Detail = detail
	} else {
		result.Detail = "Your request met some conflict."
	}

	return result
}

// 404
func NotFound(detail string) *jsonapi.ErrorObject {
	result := &jsonapi.ErrorObject{
		Title:  "Not Found",
		Status: fmt.Sprintf("%d", http.StatusNotFound),
	}

	if len(detail) > 0 {
		result.Detail = detail
	} else {
		result.Detail = "Not Found."
	}

	return result
}

// 409
func Conflict(detail string) *jsonapi.ErrorObject {
	result := &jsonapi.ErrorObject{
		Title:  "Conflict",
		Status: fmt.Sprintf("%d", http.StatusConflict),
	}

	if len(detail) > 0 {
		result.Detail = detail
	} else {
		result.Detail = "Your request met some conflict."
	}

	return result
}

func UnsupportedMediaType() *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Title:  "Unsupported Media Type",
		Detail: "Unsupported Media Type",
		Status: fmt.Sprintf("%d", http.StatusUnsupportedMediaType),
	}
}

// 500
func ServerError(err error) *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Title:  "Internal Server Error",
		Detail: "Internal Server Error",
		Status: fmt.Sprintf("%d", http.StatusInternalServerError),
	}
}
