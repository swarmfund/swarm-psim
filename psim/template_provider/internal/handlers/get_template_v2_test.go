package handlers

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"io"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/template_provider/internal/handlers/mocks"
	"gitlab.com/swarmfund/psim/psim/template_provider/internal/middlewares"
)

func TestGetTemplateV2(t *testing.T) {
	cases := []struct {
		name       string
		bucket     string
		key        string
		body       string
		statusCode int
		err        error
	}{
		{
			name:   "valid",
			key:    "template",
			bucket: "bucket",
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
			statusCode: 200,
		},
		{
			name:       "failed to download",
			key:        "template",
			bucket:     "bucket",
			err:        errors.New("failed to download"),
			statusCode: 500,
		},
		{
			name:       "invalid key",
			key:        "template",
			bucket:     "bucket",
			err:        awserr.New(s3.ErrCodeNoSuchKey, "", nil),
			statusCode: 404,
		},
		{
			name:       "invalid bucket",
			key:        "template",
			bucket:     "bucket",
			err:        awserr.New(s3.ErrCodeNoSuchBucket, "", nil),
			statusCode: 500,
		},
		{
			name:       "old template",
			key:        "template",
			bucket:     "bucket",
			body:       "body",
			statusCode: 409,
		},
	}

	downloader := &mocks.TemplateDownloader{}
	logger := logan.New()
	router := chi.NewRouter()
	router.Use(
		middlewares.Ctx(
			CtxBucket("bucket"),
			CtxDownloader(downloader),
			CtxLog(logger),
		),
	)
	router.Get("/v2/templates/{template}", GetTemplateV2)

	ts := httptest.NewServer(router)
	defer ts.Close()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockDownloadFunc := func(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) error {
				if tc.err != nil {
					return tc.err
				}
				w.WriteAt([]byte(tc.body), 0)
				return nil
			}

			downloader.On("Download",
				mock.Anything, mock.Anything).
				Return(int64(0), mockDownloadFunc).Once()
			defer downloader.AssertExpectations(t)

			resp := Client(t, ts).Do("GET", fmt.Sprintf("v2/templates/%s", tc.key), "")
			assert.Equal(t, tc.statusCode, resp.StatusCode)

		})
	}
}
