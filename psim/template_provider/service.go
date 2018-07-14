package template_provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/internal"
	"gitlab.com/swarmfund/psim/psim/template_provider/internal/handlers"
	"gitlab.com/swarmfund/psim/psim/template_provider/internal/middlewares"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/go/resources"
	"gitlab.com/tokend/horizon-connector"
)

// AccountQ interface for doorman initialization
type AccountQ interface {
	Signers(address string) ([]resources.Signer, error)
}

type Service struct {
	config     *Config
	uploader   *s3.S3
	downloader *s3manager.Downloader
	log        *logan.Entry
	info       *horizon.Info

	infoer  internal.Infoer
	doorman doorman.Doorman
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

	r.Get("/v2/templates/{template}", handlers.GetTemplateV2)
	r.Put("/v2/templates/{template}", handlers.PutTemplateV2)

	return r
}

func New(
	session *session.Session,
	log *logan.Entry,
	config *Config,
	infoer internal.Infoer,
	doorman doorman.Doorman,
) *Service {
	return &Service{
		config:     config,
		uploader:   s3.New(session),
		downloader: s3manager.NewDownloader(session),
		log:        log,
		infoer:     infoer,
		doorman:    doorman,
	}
}

func (s *Service) Run(ctx context.Context) {
	r := Router(
		s.log,
		s.uploader,
		s.downloader,
		s.config.Bucket,
		s.info,
		s.doorman,
	)

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	if err := http.ListenAndServe(addr, r); err != nil {
		s.log.WithError(err).Error("listen and serve died")
		return
	}
}
