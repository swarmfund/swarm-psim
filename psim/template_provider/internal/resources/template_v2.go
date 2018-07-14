package resources

import . "github.com/go-ozzo/ozzo-validation"

type (
	TemplateV2 struct {
		Data TemplateV2Data `json:"data"`
	}

	TemplateV2Data struct {
		Attributes TemplateV2Attributes `json:"attributes"`
	}
	TemplateV2Attributes struct {
		Subject   string `json:"subject"`
		Body      string `json:"body"`
		CreatedAt string `json:"created_at,omitempty"`
	}
)

func (p TemplateV2) Validate() error {
	return ValidateStruct(&p,
		Field(&p.Data, Required),
	)
}

func (d TemplateV2Data) Validate() error {
	return ValidateStruct(&d,
		Field(&d.Attributes, Required),
	)
}

func (a TemplateV2Attributes) Validate() error {
	return ValidateStruct(&a,
		Field(&a.Subject, Required),
		Field(&a.Body, Required),
	)
}
