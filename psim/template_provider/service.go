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
	"gitlab.com/swarmfund/go/doorman"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/template_provider/data"
)

type Service struct {
	API        TemplateAPI
	uploader   *s3.S3
	downloader *s3manager.Downloader
	log        *logan.Entry
	horizon    *horizon.Connector
}

func Router(log *logan.Entry, uploader *s3.S3, downloader *s3manager.Downloader,
	bucket string, info *horizon.Info, doorman doorman.Doorman) chi.Router {
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

func New(sess *session.Session, log *logan.Entry, api TemplateAPI, horizon *horizon.Connector) *Service {
	return &Service{
		API:        api,
		uploader:   s3.New(sess),
		downloader: s3manager.NewDownloader(sess),
		log:        log,
		horizon:    horizon,
	}
}

func (s *Service) Run(ctx context.Context) {

	info, err := s.horizon.Info()
	if err != nil {
		s.log.WithField("Failed to load", "horizon info").WithError(err).Info()
		return
	}

	r := Router(
		s.log,
		s.uploader,
		s.downloader,
		s.API.Bucket,
		info,
		doorman.New(
			s.API.SkipSignatureCheck,
			data.NewAccountQ(s.horizon),
		),
	)

	addr := fmt.Sprintf("%s:%d", s.API.Host, s.API.Port)
	if err := http.ListenAndServe(addr, r); err != nil {
		s.log.WithField("CloudAPI", "failed").WithError(err).Error()
		return
	}
}
