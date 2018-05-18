package template_provider

import (
	"context"

	"fmt"

	"net/http"

	"github.com/go-chi/chi"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/template_provider/handlers"
	"gitlab.com/swarmfund/psim/psim/template_provider/middlewares"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/horizon-connector"
)

type Service struct {
	API        Config
	uploader   *s3.S3
	downloader *s3manager.Downloader
	log        *logan.Entry
	horizon    *horizon.Connector
	info       *horizon.Info
}

func Router(
	log *logan.Entry,
	uploader *s3.S3,
	downloader *s3manager.Downloader,
	bucket string,
	info *horizon.Info,
	doorman doorman.Doorman) chi.Router {

	r := chi.NewRouter()

	r.Use(
		middlewares.ContenType("text/html"),
		ape.RecoverMiddleware(log),
		middlewares.Logger(log),
		middlewares.Ctx(
			handlers.CtxUploader(uploader),
			handlers.CtxDownloader(downloader),
			handlers.CtxLog(log),
			handlers.CtxBucket(bucket),
			handlers.CtxDoorman(doorman),
			handlers.CtxHorizonInfo(info),
		),
	)

	r.Get("/templates/{template}", handlers.GetTemplate)
	r.Put("/templates/{template}", handlers.PutTemplate)

	return r
}

func New(sess *session.Session, log *logan.Entry, api Config, info *horizon.Info, horizon *horizon.Connector) *Service {
	return &Service{
		API:        api,
		uploader:   s3.New(sess),
		downloader: s3manager.NewDownloader(sess),
		log:        log,
		horizon:    horizon,
		info:       info,
	}
}

func (s *Service) Run(ctx context.Context) {
	metric := app.Metric(ctx)
	r := Router(
		s.log,
		s.uploader,
		s.downloader,
		s.API.Bucket,
		s.info,
		doorman.New(s.API.SkipSignatureCheck, s.horizon.Accounts()),
	)

	addr := fmt.Sprintf("%s:%d", s.API.Host, s.API.Port)
	if err := http.ListenAndServe(addr, r); err != nil {
		metric.Unhealthy(err)
		s.log.WithError(err).Error("failed to listen and serve")
		return
	}
}
