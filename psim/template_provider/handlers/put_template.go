package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"io/ioutil"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/go/doorman"
)

func PutTemplate(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "template")
	bucket := Bucket(r)

	if err := Doorman(r, doorman.SignerOf(Info(r).MasterAccountID)); err != nil {
		Log(r).WithField("Failed doorman", "signature check").Error()
		RenderDoormanErr(w, err)
		return
	}

	uploader := Uploader(r)

	template, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Log(r).WithField("Can't read", "body").WithError(err).Error()
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	_, err = uploader.PutObject(&s3.PutObjectInput{
		Body:   strings.NewReader(string(template)),
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		Log(r).WithField("Failed to upload data", fmt.Sprintf("%s %s", bucket, key)).WithError(err).Error()
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
