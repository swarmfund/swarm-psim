package handlers

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func GetTemplate(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "template")
	if len(key) == 0 {
		ape.RenderErr(w, problems.BadRequest(errors.New("invalid key"))...)
		return
	}

	bucket := Bucket(r)

	file := &aws.WriteAtBuffer{}
	_, err := Downloader(r).Download(file,
		&s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
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
		Log(r).WithFields(logan.F{"bucket": bucket, "key": key}).WithError(err).Error("Failed to download")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.Write(file.Bytes())
}
