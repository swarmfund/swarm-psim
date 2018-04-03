package handlers

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/go/doorman"
)

func GetTemplate(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "template")
	bucket := Bucket(r)

	if err := Doorman(r, doorman.SignerOf(Info(r).MasterAccountID)); err != nil {
		Log(r).WithField("Failed doorman", "signature check").Error()
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
		Log(r).WithField("Failed to download file", key).WithError(err).Error()
		ape.RenderErr(w, problems.InternalError()) //todo assert errors
		return
	}

	template := file.Bytes()
	w.Write(template)
}
