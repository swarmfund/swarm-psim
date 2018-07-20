package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplateV2_Validate(t *testing.T) {
	cases := []struct {
		name      string
		errString string
		err       bool
		template  TemplateV2
	}{
		{
			name: "valid",
			template: TemplateV2{
				Data: TemplateV2Data{
					Attributes: TemplateV2Attributes{
						Body:    "body",
						Subject: "subject",
					},
				},
			},
		},
		{
			name:      "invalid",
			err:       true,
			errString: "data: (attributes: (body: cannot be blank; subject: cannot be blank.).).",
			template:  TemplateV2{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			request := TemplateV2{
				Data: tc.template.Data,
			}
			err := request.Validate()
			if tc.err && assert.Error(t, err) {
				assert.Equal(t, tc.errString, err.Error())
			}
		})
	}
}
