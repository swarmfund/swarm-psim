package handlers

import (
	"net/http"
	"strings"

	"encoding/json"

	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/template_provider/internal/resources"
	"gitlab.com/tokend/go/doorman"
)

func PutTemplateV2(w http.ResponseWriter, r *http.Request) {
	request, err := resources.NewPutTemplateRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if err := Doorman(r, doorman.SignerOf(Info(r).MasterAccountID)); err != nil {
		RenderDoormanErr(w, err)
		return
	}

	request.Data.Attributes.CreatedAt = time.Now().Format(time.RFC3339)
	body, err := json.Marshal(request)
	if err != nil {
		Log(r).WithError(err).Error("Can't marshal request")
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
		Log(r).WithFields(logan.F{"bucket": bucket, "key": request.Key}).
			WithError(err).
			Error("Failed to Upload")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
