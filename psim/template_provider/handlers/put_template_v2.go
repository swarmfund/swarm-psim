package handlers

import (
	"net/http"
	"strings"

	"encoding/json"

	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/doorman"
)

func NewTemplateV2(r *http.Request) (PutTemplateV2Req, error) {
	request := PutTemplateV2Req{
		Key: chi.URLParam(r, "template"),
	}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	request.Data.Attributes.CreatedAt = time.Now().String()
	return request, request.Validate()
}

func PutTemplateV2(w http.ResponseWriter, r *http.Request) {

	template, err := NewTemplateV2(r)

	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if err := Doorman(r, doorman.SignerOf(Info(r).MasterAccountID)); err != nil {
		RenderDoormanErr(w, err)
		return
	}

	body, err := json.Marshal(template)
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
		Key:    &template.Key,
	})

	if err != nil {
		Log(r).WithFields(logan.F{"bucket": bucket, "key": template.Key}).
			WithError(err).
			Error("Failed to Upload")
		ape.RenderErr(w, problems.InternalError())
	}

	w.WriteHeader(http.StatusNoContent)
}
