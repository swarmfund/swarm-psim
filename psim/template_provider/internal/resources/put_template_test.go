package resources

import (
	"reflect"
	"testing"

	"gitlab.com/distributed_lab/ape/apeutil"
)

func TestNewPutTemplateRequest(t *testing.T) {
	cases := []struct {
		name     string
		key      string
		body     string
		err      bool
		expected PutTemplateRequest
	}{
		{
			name: "valid request",
			key:  "key",
			err:  false,
			body: `{
						"data":
						{
			         		"attributes":
							{
			             		"body": "body",
			             		"subject": "subject"
							}
						}
			  		}`,
			expected: PutTemplateRequest{
				Key: "key",
				Data: TemplateV2Data{
					Attributes: TemplateV2Attributes{
						Body:    "body",
						Subject: "subject",
					},
				},
			},
		},
		{
			name: "no key",
			key:  "",
			err:  true,
			body: `{
						"data":
						{
			         		"attributes":
							{
			             		"body": "body",
			             		"subject": "subject"
							}
						}
			  		}`,

			expected: PutTemplateRequest{
				Key: "",
				Data: TemplateV2Data{
					Attributes: TemplateV2Attributes{
						Body:    "body",
						Subject: "subject",
					},
				},
			},
		},
		{
			name: "no body",
			key:  "key",
			err:  true,
			body: `{
						"data":
						{
			         		"attributes":
							{
			             		"subject": "subject"
							}
						}
			  		}`,

			expected: PutTemplateRequest{
				Key: "key",
				Data: TemplateV2Data{
					Attributes: TemplateV2Attributes{
						Subject: "subject",
					},
				},
			},
		},
		{
			name: "no subject",
			key:  "key",
			err:  true,
			body: `{
						"data":
						{
			         		"attributes":
							{
			             		"body": "body",
							}
						}
			  		}`,
			expected: PutTemplateRequest{
				Key: "",
				Data: TemplateV2Data{
					Attributes: TemplateV2Attributes{
						Body: "body",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := apeutil.RequestWithURLParams([]byte(tc.body), map[string]string{
				"template": tc.key,
			})
			got, err := NewPutTemplateRequest(r)
			if err != nil && !tc.err {
				t.Fatalf("expected nil error got %s", err)
			}
			if err == nil && tc.err {
				t.Fatalf("expected error got nil")
			}
			if err == nil && !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("expected %#v got #%v", tc.expected, got)
			}
		})
	}
}
