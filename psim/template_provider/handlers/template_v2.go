package handlers

import (
	. "github.com/go-ozzo/ozzo-validation"
)

type (
	GetTemplateV2Req struct {
		Key string `json:"-"`
	}
	PutTemplateV2Req struct {
		Key      string     `json:"-"`
		Template TemplateV2 `json:"template"`
	}

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

func (p PutTemplateV2Req) Validate() error {
	return ValidateStruct(&p,
		Field(&p.Key, Required),
		Field(&p.Template, Required),
	)
}

func (p TemplateV2) Validate() error {
	return ValidateStruct(&p,
		Field(&p.Data, Required),
	)
}

func (r GetTemplateV2Req) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Key, Required),
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
