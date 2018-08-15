package handlers

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3"

	"gitlab.com/swarmfund/psim/psim/template_provider/internal/middlewares"
	"gitlab.com/swarmfund/psim/psim/template_provider/internal/resources"
	"gitlab.com/tokend/go/doorman"
)

func Router(
	log *logan.Entry,
	uploader resources.TemplateUploader,
	downloader resources.TemplateDownloader,
	bucket string,
	info *regources.Info,
	doorman doorman.Doorman) chi.Router {

	r := chi.NewRouter()

	r.Use(
		middlewares.ContenType("text/html"),
		ape.RecoverMiddleware(log),
		middlewares.Logger(log),
		middlewares.Ctx(
			CtxUploader(uploader),
			CtxDownloader(downloader),
			CtxLog(log),
			CtxBucket(bucket),
			CtxDoorman(doorman),
			CtxHorizonInfo(info),
		),
	)

	r.Get("/templates/{template}", GetTemplate)
	r.Put("/templates/{template}", PutTemplate)

	r.Get("/v2/templates/{template}", GetTemplateV2)
	r.Put("/v2/templates/{template}", PutTemplateV2)

	return r
}
