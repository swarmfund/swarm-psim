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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/template_provider/internal/handlers/mocks"
	"gitlab.com/swarmfund/psim/psim/template_provider/internal/middlewares"
)

func TestGetTemplate(t *testing.T) {
	cases := []struct {
		name       string
		bucket     string
		key        string
		actualKey  string
		actual     string
		expected   string
		statusCode int
		err        bool
	}{
		{
			name:       "valid",
			key:        "template",
			actualKey:  "template",
			bucket:     "bucket",
			actual:     "This is the contents of template",
			expected:   "This is the contents of template",
			statusCode: 200,
		},
		{
			name:       "Failed to download",
			key:        "template",
			actualKey:  "template",
			bucket:     "bucket",
			actual:     "This issdkfj the",
			expected:   "This is the",
			err:        true,
			statusCode: 500,
		},
		{
			name:       "invalid key",
			key:        "notemplate",
			actualKey:  "template",
			bucket:     "bucket",
			err:        true,
			statusCode: 404,
		},
		{
			name:       "invalid bucket",
			key:        "template",
			bucket:     "bukit",
			err:        true,
			statusCode: 500,
		},
	}

	downloader := &mocks.TemplateDownloader{}
	logger := logan.New()
	router := chi.NewRouter()
	router.Use(
		middlewares.ContenType("text/html"),
		middlewares.Ctx(
			CtxBucket("bucket"),
			CtxDownloader(downloader),
			CtxLog(logger),
		),
	)
	router.Get("/templates/{template}", GetTemplate)

	ts := httptest.NewServer(router)
	defer ts.Close()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			mockDownloadFunc := func(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) error {

				if tc.actual != tc.expected {
					return errors.New("Failed to download")
				}

				if tc.bucket != "bucket" {
					return awserr.New(s3.ErrCodeNoSuchBucket, "No such bucket", nil)
				}

				if tc.key != tc.actualKey {
					return awserr.New(s3.ErrCodeNoSuchKey, "", nil)
				}

				w.WriteAt([]byte(tc.actual), 0)
				return nil
			}
			downloader.On("Download",
				mock.Anything, mock.Anything).
				Return(int64(0), mockDownloadFunc).Once()
			defer downloader.AssertExpectations(t)

			resp := Client(t, ts).Do("GET", fmt.Sprintf("templates/%s", tc.key), "")
			assert.Equal(t, tc.statusCode, resp.StatusCode)

		})
	}
}
