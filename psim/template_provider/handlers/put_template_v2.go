package handlers

import (
	"net/http"
	"strings"

	"encoding/json"

	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/doorman"
)

type (
	PutTemplateV2 struct {
		Key  string
		Data TemplateV2Data `json:"data"`
	}
	TemplateV2Data struct {
		Type       string               `json:"type"`
		Attributes TemplateV2Attributes `json:"attributes"`
	}
	TemplateV2Attributes struct {
		Subject   string `json:"subject"`
		Body      string `json:"body"`
		CreatedAt string `json:"created_at,omitempty"`
	}
)

func NewPutTemplateV2Req(r *http.Request) (PutTemplateV2, error) {
	request := PutTemplateV2{
		Key: chi.URLParam(r, "template"),
	}
	if len(request.Key) == 0 {
		return request, errors.New("invalid key")
	}
	if err := json.NewDecoder(r.Body).Decode(&request.Data); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	request.Data.Attributes.CreatedAt = time.Now().String()
	return request, request.Validate()
}

func (r PutTemplateV2) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Key, Required),
		Field(&r.Data, Required),
	)
}
func (d TemplateV2Data) Validate() error {
	return ValidateStruct(&d,
		Field(&d.Type, Required),
		Field(&d.Attributes, Required),
	)
}

func (a TemplateV2Attributes) Validate() error {
	return ValidateStruct(&a,
		Field(&a.Subject, Required),
		Field(&a.Body, Required),
	)
}

func PutTemplateWithSubject(w http.ResponseWriter, r *http.Request) {

	request, err := NewPutTemplateV2Req(r)

	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if err := Doorman(r, doorman.SignerOf(Info(r).MasterAccountID)); err != nil {
		RenderDoormanErr(w, err)
		return
	}

	body, err := json.Marshal(request.Data.Attributes)
	if err != nil {
		Log(r).WithError(err).Error("Can't unmarshal request")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	bucket := Bucket(r)
	uploader := Uploader(r)
	_, err = uploader.PutObject(&s3.PutObjectInput{
		Body:   strings.NewReader(string(body)),
		Bucket: &bucket,
		Key:    &request.Key,
	})

	if err != nil {
		Log(r).WithFields(logan.F{"bucket": bucket, "key": request.Key}).WithError(err).Error("Failed to download")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
