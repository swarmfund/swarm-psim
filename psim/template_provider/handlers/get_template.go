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
	"gitlab.com/tokend/go/doorman"
)

func GetTemplate(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "template")
	if len(key) == 0 {
		ape.RenderErr(w, problems.BadRequest(errors.New("invalid key"))...)
		return
	}

	bucket := Bucket(r)

	if err := Doorman(r, doorman.SignerOf(Info(r).MasterAccountID)); err != nil {
		RenderDoormanErr(w, err)
		return
	}

	downloader := Downloader(r)

	file := &aws.WriteAtBuffer{}
	_, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				Log(r).WithFields(logan.F{"bucket": bucket}).WithError(err).Error("No such bucket")
				ape.RenderErr(w, problems.InternalError())
			case s3.ErrCodeNoSuchKey:
				ape.RenderErr(w, problems.NotFound())
			}
			return
		}
		Log(r).WithFields(logan.F{"bucket": bucket, "key": key}).WithError(err).Error("Failed to download")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	template := file.Bytes()
	w.Write(template)
}
