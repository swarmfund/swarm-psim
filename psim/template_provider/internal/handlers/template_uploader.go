package handlers

import (
	"github.com/aws/aws-sdk-go/service/s3"
)

//go:generate mockery -case underscore -name TemplateUploader

type TemplateUploader interface {
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
}
