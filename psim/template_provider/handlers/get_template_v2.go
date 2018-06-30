package handlers

import (
	"net/http"

	"encoding/json"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi"
	"github.com/google/jsonapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/doorman"
)

func NewGetTemplateV2(r *http.Request) (GetTemplateV2Req, error) {
	request := GetTemplateV2Req{
		Key: chi.URLParam(r, "template"),
	}
	return request, request.Validate()
}

func GetTemplateV2(w http.ResponseWriter, r *http.Request) {
	request, err := NewGetTemplateV2(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if err := Doorman(r, doorman.SignerOf(Info(r).MasterAccountID)); err != nil {
		RenderDoormanErr(w, err)
		return
	}

	bucket := Bucket(r)
	downloader := Downloader(r)

	file := &aws.WriteAtBuffer{}
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &request.Key,
		})
	if err != nil {
		cause := errors.Cause(err)
		if aerr, ok := cause.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				Log(r).WithFields(logan.F{"bucket": bucket}).WithError(err).Error("No such bucket")
				ape.RenderErr(w, problems.InternalError())
				return
			case s3.ErrCodeNoSuchKey:
				ape.RenderErr(w, problems.NotFound())
				return
			}
		}
		Log(r).WithFields(logan.F{"bucket": bucket, "key": request.Key}).WithError(err).Error("Failed to download")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	raw := file.Bytes()
	var template TemplateV2
	err = json.Unmarshal(raw, &template)
	if err != nil {
		Log(r).WithFields(logan.F{"bucket": bucket, "key": request.Key}).WithError(err).Error("Failed to download")
		ape.RenderErr(w, &jsonapi.ErrorObject{
			Title:  http.StatusText(http.StatusSeeOther),
			Status: fmt.Sprintf("%d", http.StatusSeeOther),
			Detail: "Resource not updated. Use /template endpoint",
		})
		return
	}
	err = template.Validate()
	if err != nil {
		Log(r).WithFields(logan.F{"bucket": bucket, "key": request.Key}).
			WithError(err).
			Error("Incorrect template format")
		ape.RenderErr(w, problems.Conflict())
		return
	}

	json.NewEncoder(w).Encode(template)

}
