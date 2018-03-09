package server

import (
	"github.com/Sirupsen/logrus"
	"gitlab.com/distributed_lab/notificator-server/conf"
	"gitlab.com/distributed_lab/notificator-server/types"
	"gitlab.com/distributed_lab/notificator-server/workers"
)

type WorkerFunc func(request types.Request, cfg conf.Config) bool

var workersMap = map[string]WorkerFunc{
	"dummy": workers.Dummy,
	//DEPRECATED don't use "email" worker!
	"email":          workers.MandrillEmail,
	"mandrill_email": workers.MandrillEmail,
	"mailgun":        workers.MailgunEmail,
	"postage":        workers.PostageEmail,
	//DEPRECATED don't use "sms" worker!
	"sms":    workers.SMS,
	"twilio": workers.SMS,
}

type Worker struct {
	tasks   <-chan types.Request
	results chan<- TaskResult
	log     *logrus.Entry
}

func NewWorker(tasks <-chan types.Request, results chan<- TaskResult, log *logrus.Entry) *Worker {
	return &Worker{
		tasks:   tasks,
		results: results,
		log:     log.WithField("service", "worker"),
	}
}

func (w *Worker) Run(cfg conf.Config) {
	var success bool

	for request := range w.tasks {
		w.log.WithField("request", request.ID).Info("processing")

		success = w.executeWorker(cfg, request)

		w.log.WithField("success", success).Info("done")
		w.results <- TaskResult{
			Request: request,
			Success: success,
		}
	}
}

func (w *Worker) executeWorker(cfg conf.Config, request types.Request) bool {
	entry := w.log.WithField("request", request.ID)
	requestsConf := cfg.Requests()

	requestType, ok := requestsConf.Get(request.Type)
	if !ok {
		entry.Error("unknown request type")
		return false
	}
	workman, present := workersMap[requestType.Worker]
	if !present {
		entry.Error("unknown worker type")
	}

	return workman(request, cfg)
}
