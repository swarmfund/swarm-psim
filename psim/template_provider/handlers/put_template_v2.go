package handlers

import (
	"net/http"
	"strings"

	"encoding/json"

	"time"

	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/doorman"
)

func NewPutTemplateV2(r *http.Request) (PutTemplateV2Req, error) {
	request := PutTemplateV2Req{
		Bucket: BucketKey{chi.URLParam(r, "template")},
	}
	if err := json.NewDecoder(r.Body).Decode(&request.Template); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	request.Template.Data.Attributes.CreatedAt = time.Now().String()
	return request, request.Validate()
}

func PutTemplateV2(w http.ResponseWriter, r *http.Request) {

	template, err := NewPutTemplateV2(r)

	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if err := Doorman(r, doorman.SignerOf(Info(r).MasterAccountID)); err != nil {
		RenderDoormanErr(w, err)
		return
	}

	body, err := json.Marshal(template.Template)
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
		Key:    &template.Bucket.Key,
	})

	if err != nil {
		Log(r).WithFields(logan.F{"bucket": bucket, "key": template.Bucket.Key}).WithError(err).Error("Failed to download")
		ape.RenderErr(w, &jsonapi.ErrorObject{
			Title:  http.StatusText(http.StatusSeeOther),
			Status: fmt.Sprintf("%d", http.StatusSeeOther),
			Detail: "Resource not updated. Use old API to access the resource",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
