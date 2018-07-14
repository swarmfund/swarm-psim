package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/template_provider/internal/resources"
	"gitlab.com/tokend/go/doorman"
)

func GetTemplateV2(w http.ResponseWriter, r *http.Request) {
	request, err := resources.NewGetTemplateRequest(r)
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
		Log(r).WithFields(logan.F{"bucket": bucket, "key": request.Key}).
			WithError(err).
			Error("Failed to unmarshal")
		ape.RenderErr(w, problems.Conflict())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(template)

}
