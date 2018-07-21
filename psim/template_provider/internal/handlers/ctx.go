package handlers

import (
	"context"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/horizon-connector"
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

func CtxUploader(uploader TemplateUploader) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, CtxUploaderKey, uploader)
	}
}

func Uploader(r *http.Request) TemplateUploader {
	return r.Context().Value(CtxUploaderKey).(TemplateUploader)
}

func CtxDownloader(downloader TemplateDownloader) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, CtxDownloaderKey, downloader)
	}
}

func Downloader(r *http.Request) TemplateDownloader {
	return r.Context().Value(CtxDownloaderKey).(TemplateDownloader)
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
