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
	PutTemplateWSubjectRequest struct {
		Key  string                  `json:"key"`
		Data PutTemplateWSubjectData `json:"data"`
	}
	PutTemplateWSubjectData struct {
		Subject   string `json:"subject"`
		Body      string `json:"body"`
		CreatedAt string `json:"created_at"`
	}
)

func NewPutTemplateWSubjectRequest(r *http.Request) (PutTemplateWSubjectRequest, error) {
	request := PutTemplateWSubjectRequest{
		Key: chi.URLParam(r, "template"),
	}
	if len(request.Key) == 0 {
		return request, errors.New("invalid key")
	}
	if err := json.NewDecoder(r.Body).Decode(&request.Data); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	request.Data.CreatedAt = time.Now().String()
	return request, request.Validate()
}

func (r PutTemplateWSubjectRequest) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Key, Required),
		Field(&r.Data, Required),
	)
}
func (d PutTemplateWSubjectData) Validate() error {
	return ValidateStruct(&d,
		Field(&d.Body, Required),
		Field(&d.Subject, Required),
	)
}

func PutTemplateWithSubject(w http.ResponseWriter, r *http.Request) {

	request, err := NewPutTemplateWSubjectRequest(r)

	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if err := Doorman(r, doorman.SignerOf(Info(r).MasterAccountID)); err != nil {
		RenderDoormanErr(w, err)
		return
	}

	body, err := json.Marshal(request.Data)
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
