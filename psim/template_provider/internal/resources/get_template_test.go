package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/ape/apeutil"
)

func TestNewGetTemplateRequest(t *testing.T) {
	cases := []struct {
		name      string
		key       string
		err       bool
		errString string
		expected  GetTemplateRequest
	}{
		{
			name: "valid request",
			key:  "key",
			err:  false,
			expected: GetTemplateRequest{
				Key: "key",
			},
		},
		{
			name:      "invalid request",
			key:       "",
			err:       true,
			errString: "-: cannot be blank.",
			expected: GetTemplateRequest{
				Key: "",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := apeutil.RequestWithURLParams(nil, map[string]string{
				"template": tc.key,
			})
			got, err := NewGetTemplateRequest(r)
			if tc.err && assert.Error(t, err) {
				assert.Equal(t, tc.errString, err.Error())
			}
			assert.Equal(t, tc.expected, got)
		})
	}
}
