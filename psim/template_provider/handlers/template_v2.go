package handlers

import (
	. "github.com/go-ozzo/ozzo-validation"
)

type (
	BucketKey struct {
		Key string `json:"-"`
	}
	TemplateV2 struct {
		Data TemplateV2Data `json:"data"`
	}
	PutTemplateV2Req struct {
		Template TemplateV2
		Bucket   BucketKey
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

func (p PutTemplateV2Req) Validate() error {
	return ValidateStruct(&p,
		Field(&p.Bucket, Required),
		Field(&p.Template, Required),
	)
}

func (r BucketKey) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Key, Required),
	)
}

func (t TemplateV2) Validate() error {
	return ValidateStruct(&t,
		Field(&t.Data, Required),
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
