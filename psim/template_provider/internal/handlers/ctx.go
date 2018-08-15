package handlers

import (
	"context"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/template_provider/internal/resources"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/regources"
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

func CtxUploader(uploader resources.TemplateUploader) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, CtxUploaderKey, uploader)
	}
}

func Uploader(r *http.Request) resources.TemplateUploader {
	return r.Context().Value(CtxUploaderKey).(resources.TemplateUploader)
}

func CtxDownloader(downloader resources.TemplateDownloader) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, CtxDownloaderKey, downloader)
	}
}

func Downloader(r *http.Request) resources.TemplateDownloader {
	return r.Context().Value(CtxDownloaderKey).(resources.TemplateDownloader)
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

func CtxHorizonInfo(i *regources.Info) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, horizonInfoCtxKey, i)
	}
}

func Info(r *http.Request) *regources.Info {
	return r.Context().Value(horizonInfoCtxKey).(*regources.Info)
}
