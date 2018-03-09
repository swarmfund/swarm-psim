package server

import (
	"github.com/Sirupsen/logrus"
	"gitlab.com/distributed_lab/notificator-server/conf"
	"gitlab.com/distributed_lab/notificator-server/types"
)

type WorkerHerder struct {
	log *logrus.Entry
}

func NewWorkerHerder(config conf.Config) *WorkerHerder {
	return &WorkerHerder{
		log: config.Log().WithField("service", "worker_herder"),
	}
}

func (h *WorkerHerder) Init(_ conf.Config) {

}

func (h *WorkerHerder) Run(cfg conf.Config) {
	workersCount := 10
	tasks := make(chan types.Request)
	results := make(chan TaskResult)
	consumed := make(chan types.Request)
	workers := make([]*Worker, workersCount)
	for i := range workers {
		workers[i] = NewWorker(tasks, results, h.log)
		go workers[i].Run(cfg)
	}

	resultsConsumer := NewTaskResultConsumer(results, consumed, h.log)
	go resultsConsumer.Run()

	taskProvider := NewTaskProvider(tasks, consumed, h.log)
	taskProvider.Run()
}
