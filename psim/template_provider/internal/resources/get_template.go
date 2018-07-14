package resources

import (
	"net/http"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
)

type GetTemplateRequest struct {
	Key string `json:"-"`
}

func NewGetTemplateRequest(r *http.Request) (GetTemplateRequest, error) {
	request := GetTemplateRequest{
		Key: chi.URLParam(r, "template"),
	}
	return request, request.Validate()
}

func (r GetTemplateRequest) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Key, Required),
	)
}
