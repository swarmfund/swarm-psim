package handlers

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/tokend/go/keypair"
	"gitlab.com/tokend/go/resources"
)

func TestPutTemplate(t *testing.T) {
	cases := []struct {
		name          string
		key           string
		body          string
		notAuthorized bool
		signer        keypair.KP
		statusCode    int
	}{
		{
			name:       "valid",
			key:        "template",
			body:       "body",
			statusCode: 204,
		},
		{
			name:       "no body",
			key:        "template",
			body:       "",
			statusCode: 500,
		},
		{
			name:          "not authorized",
			key:           "template",
			body:          "body",
			statusCode:    401,
			notAuthorized: true,
		},
	}

	r, uploader, _, signer, accountQ := testRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var caseSigner keypair.KP
			if tc.notAuthorized {
				caseSigner, _ = keypair.Random()
			} else {
				caseSigner = signer
			}

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
			if !tc.notAuthorized {
				defer uploader.AssertExpectations(t)
			}

			sign := func(account string) []resources.Signer {
				return []resources.Signer{
					{
						AccountID: signer.Address(),
						Weight:    1,
					},
				}
			}
			accountQ.On("Signers", mock.MatchedBy(func(address string) bool {
				return len(address) != 0
			})).Return(sign, nil)

			resp := Client(t, ts).Signer(caseSigner).Do("PUT", fmt.Sprintf("templates/%s", tc.key), tc.body)
			assert.Equal(t, tc.statusCode, resp.StatusCode)
		})
	}
}
