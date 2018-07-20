package resources

import "testing"

func TestTemplateV2_Validate(t *testing.T) {
	cases := []struct {
		name     string
		err      bool
		template TemplateV2
	}{
		{
			name: "valid",
			err:  false,
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
			name:     "invalid",
			err:      true,
			template: TemplateV2{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			request := TemplateV2{
				Data: tc.template.Data,
			}
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

func TestTemplateV2Data_Validate(t *testing.T) {
	cases := []struct {
		name string
		err  bool
		data TemplateV2Data
	}{
		{
			name: "valid",
			err:  false,
			data: TemplateV2Data{
				Attributes: TemplateV2Attributes{
					Subject: "subject",
					Body:    "body",
				},
			},
		},
		{
			name: "invalid",
			err:  true,
			data: TemplateV2Data{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			request := TemplateV2Data{
				Attributes: tc.data.Attributes,
			}
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

func TestTemplateV2Attributes_Validate(t *testing.T) {
	cases := []struct {
		name       string
		attributes TemplateV2Attributes
		err        bool
	}{
		{
			name: "valid",
			attributes: TemplateV2Attributes{
				Subject: "subject",
				Body:    "body",
			},
			err: false,
		},
		{
			name: "no body",
			attributes: TemplateV2Attributes{
				Subject: "subject",
			},
			err: true,
		},
		{
			name: "no subject",
			attributes: TemplateV2Attributes{
				Body: "body",
			},
			err: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			request := TemplateV2Attributes{
				Subject: tc.attributes.Subject,
				Body:    tc.attributes.Body,
			}
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
