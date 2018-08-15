package template_provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/internal"
	"gitlab.com/swarmfund/psim/psim/template_provider/internal/handlers"
	res "gitlab.com/swarmfund/psim/psim/template_provider/internal/resources"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/go/resources"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/regources"
)

// AccountQ interface for doorman initialization
type AccountQ interface {
	Signers(address string) ([]resources.Signer, error)
}

type Service struct {
	config     *Config
	uploader   res.TemplateUploader
	downloader res.TemplateDownloader
	log        *logan.Entry
	info       *regources.Info

	infoer  internal.Infoer
	doorman doorman.Doorman
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
	r := handlers.Router(
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
