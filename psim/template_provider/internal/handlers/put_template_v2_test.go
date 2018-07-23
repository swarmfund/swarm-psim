package handlers

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/template_provider/internal/handlers/mocks"
	"gitlab.com/swarmfund/psim/psim/template_provider/internal/middlewares"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/go/keypair"
	"gitlab.com/tokend/horizon-connector"
)

func TestPutTemplateV2(t *testing.T) {
	cases := []struct {
		name       string
		key        string
		body       string
		statusCode int
		err        error
	}{
		{
			name: "valid",
			key:  "template",
			body: `{
							"data":
							{
								"attributes":
								{
									"body": "body of template",
									"subject": "subject"
								}
							}
						}`,
			statusCode: 204,
		},
		{
			name: "failed to upload",
			key:  "template",
			body: `{
							"data":
							{
								"attributes":
								{
									"body": "body of template",
									"subject": "subject"
								}
							}
						}`,
			err:        errors.New("failed to upload"),
			statusCode: 500,
		},
	}

	signer, err := keypair.Random()
	if err != nil {
		t.Fatal(err)
	}

	accountQ := mocks.AccountQ{}
	doormanM := doorman.New(
		false, &accountQ,
	)
	uploader := &mocks.TemplateUploader{}
	logger := logan.New()
	info := &horizon.Info{
		MasterAccountID: signer.Address(),
	}

	router := chi.NewRouter()
	router.Use(
		middlewares.Ctx(
			CtxBucket("bucket"),
			CtxUploader(uploader),
			CtxLog(logger),
			CtxDoorman(doormanM),
			CtxHorizonInfo(info),
		),
	)
	router.Put("/v2/templates/{template}", PutTemplateV2)

	ts := httptest.NewServer(router)
	defer ts.Close()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockUploadFunc := func(input *s3.PutObjectInput) error {
				if tc.err != nil {
					return tc.err
				}
				return nil
			}
			uploader.On("PutObject",
				mock.MatchedBy(func(input *s3.PutObjectInput) bool {
					return input.Key != nil && input.Bucket != nil && input.Body != nil
				})).
				Return(nil, mockUploadFunc).Once()
			defer uploader.AssertExpectations(t)

			resp := Client(t, ts).Signer(signer).Do("PUT", fmt.Sprintf("v2/templates/%s", tc.key), tc.body)
			assert.Equal(t, tc.statusCode, resp.StatusCode)

		})
	}
}
