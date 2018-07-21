package handlers

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//go:generate mockery -case underscore -name TemplateDownloader

type TemplateDownloader interface {
	Download(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error)
	DownloadWithContext(ctx aws.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error)
	DownloadWithIterator(ctx aws.Context, iter s3manager.BatchDownloadIterator, opts ...func(*s3manager.Downloader)) error
}
