package handlers

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"io"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetTemplate(t *testing.T) {
	cases := []struct {
		name       string
		key        string
		body       string
		statusCode int
		err        error
	}{
		{
			name:       "valid",
			key:        "template",
			body:       "This is the contents of template",
			statusCode: 200,
		},
		{
			name:       "failed to download",
			key:        "template",
			err:        errors.New("Failed to download"),
			statusCode: 500,
		},
		{
			name:       "invalid key",
			key:        "template",
			err:        awserr.New(s3.ErrCodeNoSuchKey, "", nil),
			statusCode: 404,
		},
		{
			name:       "invalid bucket",
			key:        "template",
			err:        awserr.New(s3.ErrCodeNoSuchBucket, "", nil),
			statusCode: 500,
		},
	}

	r, _, downloader, _, _ := testRouter()

	ts := httptest.NewServer(r)
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
				mock.MatchedBy(func(w io.WriterAt) bool {
					return w != nil
				}),
				mock.MatchedBy(func(input *s3.GetObjectInput) bool {
					return input.Bucket != nil && input.Key != nil
				})).
				Return(int64(0), mockDownloadFunc).Once()
			defer downloader.AssertExpectations(t)

			resp := Client(t, ts).Do("GET", fmt.Sprintf("templates/%s", tc.key), "")
			assert.Equal(t, tc.statusCode, resp.StatusCode)

		})
	}
}
