package resources

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
)

type PutTemplateRequest struct {
	Key  string         `json:"-"`
	Data TemplateV2Data `json:"data"`
}

func NewPutTemplateRequest(r *http.Request) (PutTemplateRequest, error) {
	request := PutTemplateRequest{
		Key: chi.URLParam(r, "template"),
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.Validate()
}

func (p PutTemplateRequest) Validate() error {
	return ValidateStruct(&p,
		Field(&p.Key, Required),
		Field(&p.Data, Required),
	)
}
