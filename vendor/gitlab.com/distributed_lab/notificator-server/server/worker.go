package server

import (
	"github.com/Sirupsen/logrus"
	"gitlab.com/distributed_lab/notificator-server/log"
	"gitlab.com/distributed_lab/notificator-server/types"
	"gitlab.com/distributed_lab/notificator-server/workers"
)

type WorkerFunc func(request types.Request) bool

var workersMap = map[string]WorkerFunc{
	"dummy":   workers.Dummy,
	"email":   workers.MandrillEmail,
	"mailgun": workers.MailgunEmail,
	"postage": workers.PostageEmail,
}

type Worker struct {
	tasks   <-chan types.Request
	results chan<- TaskResult
	log     *logrus.Entry
}

func NewWorker(tasks <-chan types.Request, results chan<- TaskResult) *Worker {
	return &Worker{
		tasks:   tasks,
		results: results,
		log:     log.WithField("service", "worker"),
	}
}

func (w *Worker) Run() {
	var success bool

	for request := range w.tasks {
		entry := w.log.WithField("request", request.ID)
		entry.Info("processing")

		success = w.executeWorker(request)

		entry.WithField("success", success).Info("done")
		w.results <- TaskResult{
			Request: request,
			Success: success,
		}
	}
}

func (w *Worker) executeWorker(request types.Request) bool {
	entry := w.log.WithField("request", request.ID)
	requestsConf := GetRequestsConf()

	requestType, ok := requestsConf.Get(request.Type)
	if !ok {
		entry.Error("unknown request type")
		return false
	}
	workman, present := workersMap[requestType.Worker]
	if !present {
		entry.Error("unknown worker type")
	}

	return workman(request)
}
