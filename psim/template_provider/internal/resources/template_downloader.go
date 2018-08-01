package resources

import (
	"io"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//go:generate mockery -case underscore -name TemplateDownloader

type TemplateDownloader interface {
	Download(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error)
}
