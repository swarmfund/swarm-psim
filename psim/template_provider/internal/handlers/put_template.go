package handlers

import (
	"net/http"
	"strings"

	"io/ioutil"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/doorman"
)

func PutTemplate(w http.ResponseWriter, r *http.Request) {
	bucket := Bucket(r)

	//TODO ADD VALIDATION KEY
	key := chi.URLParam(r, "template")
	if len(key) == 0 {
		ape.RenderErr(w, problems.BadRequest(errors.New("invalid key"))...)
		return
	}
	if err := Doorman(r, doorman.SignerOf(Info(r).MasterAccountID)); err != nil {
		RenderDoormanErr(w, err)
		return
	}

	uploader := Uploader(r)

	template, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	_, err = uploader.PutObject(&s3.PutObjectInput{
		Body:   strings.NewReader(string(template)),
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		Log(r).WithFields(logan.F{"bucket": bucket, "key": key}).WithError(err).Error("Failed to download")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
