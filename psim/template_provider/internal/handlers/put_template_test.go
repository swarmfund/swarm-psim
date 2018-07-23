package handlers

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/logan/v3"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"gitlab.com/swarmfund/psim/psim/template_provider/internal/handlers/mocks"
	"gitlab.com/swarmfund/psim/psim/template_provider/internal/middlewares"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/go/keypair"
	"gitlab.com/tokend/horizon-connector"
)

func TestPutTemplate(t *testing.T) {
	cases := []struct {
		name       string
		bucket     string
		key        string
		actualKey  string
		body       string
		statusCode int
	}{
		{
			name:       "valid",
			key:        "template",
			actualKey:  "template",
			bucket:     "bucket",
			body:       "body",
			statusCode: 204,
		},
		{
			name:       "no body",
			key:        "template",
			bucket:     "bucket",
			body:       "",
			statusCode: 500,
		},
	}

	bucket := "bucket"
	uploader := &mocks.TemplateUploader{}

	signer, err := keypair.Random()
	if err != nil {
		t.Fatal(err)
	}

	accountQ := mocks.AccountQ{}
	doormanM := doorman.New(
		true, &accountQ,
	)
	info := &horizon.Info{
		MasterAccountID: signer.Address(),
	}

	logger := logan.New()

	router := chi.NewRouter()
	router.Use(
		middlewares.ContenType("text/html"),
		middlewares.Ctx(
			CtxBucket(bucket),
			CtxUploader(uploader),
			CtxDoorman(doormanM),
			CtxHorizonInfo(info),
			CtxLog(logger),
		),
	)
	router.Put("/templates/{template}", PutTemplate)

	ts := httptest.NewServer(router)
	defer ts.Close()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			mockUploadFunc := func(input *s3.PutObjectInput) error {
				if len(tc.body) == 0 {
					return errors.New("Empty body")
				}
				return nil
			}
			uploader.On("PutObject",
				mock.MatchedBy(func(input *s3.PutObjectInput) bool {
					return input.Key != nil && input.Bucket != nil && input.Body != nil
				})).
				Return(nil, mockUploadFunc).Once()
			defer uploader.AssertExpectations(t)

			resp := Client(t, ts).Signer(signer).Do("PUT", fmt.Sprintf("templates/%s", tc.key), tc.body)
			assert.Equal(t, tc.statusCode, resp.StatusCode)
		})
	}
}
