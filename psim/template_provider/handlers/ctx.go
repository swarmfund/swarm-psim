package handlers

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/go/doorman"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

type ctxKey int

const (
	CtxUploaderKey ctxKey = iota
	horizonInfoCtxKey
	CtxDownloaderKey
	DoormanCtxKey
	CtxBucketKey
	CtxLogKey
)

func CtxUploader(uploader *s3.S3) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, CtxUploaderKey, uploader)
	}
}

func Uploader(r *http.Request) *s3.S3 {
	return r.Context().Value(CtxUploaderKey).(*s3.S3)
}

func CtxDownloader(uploader *s3manager.Downloader) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, CtxDownloaderKey, uploader)
	}
}

func Downloader(r *http.Request) *s3manager.Downloader {
	return r.Context().Value(CtxDownloaderKey).(*s3manager.Downloader)
}

func CtxLog(log *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, CtxLogKey, log)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(CtxLogKey).(*logan.Entry)
}

func CtxBucket(bucket string) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, CtxBucketKey, bucket)
	}
}

func Bucket(r *http.Request) string {
	return r.Context().Value(CtxBucketKey).(string)
}

func CtxDoorman(d doorman.Doorman) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, DoormanCtxKey, d)
	}
}

func Doorman(r *http.Request, constraints ...doorman.SignerConstraint) error {
	d := r.Context().Value(DoormanCtxKey).(doorman.Doorman)
	return d.Check(r, constraints...)
}

func CtxHorizonInfo(i *horizon.Info) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, horizonInfoCtxKey, i)
	}
}

func Info(r *http.Request) *horizon.Info {
	return r.Context().Value(horizonInfoCtxKey).(*horizon.Info)
}
