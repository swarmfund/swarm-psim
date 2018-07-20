package resources

import (
	"testing"

	"reflect"

	"gitlab.com/distributed_lab/ape/apeutil"
)

func TestGetTemplateRequest_Validate(t *testing.T) {
	cases := []struct {
		name string
		key  string
		err  bool
	}{
		{
			name: "valid",
			key:  "key",
			err:  false,
		},
		{
			name: "invalid",
			key:  "",
			err:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			request := GetTemplateRequest{Key: tc.key}
			err := request.Validate()
			if err != nil && !tc.err {
				t.Fatalf("expected nil error got %s", err)
			}
			if err == nil && tc.err {
				t.Fatalf("expected error got nil")
			}
		})
	}

}

func TestNewGetTemplateRequest(t *testing.T) {
	cases := []struct {
		name     string
		key      string
		err      bool
		expected GetTemplateRequest
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
			name: "invalid request",
			key:  "",
			err:  true,
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
