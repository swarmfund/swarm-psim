package server

import (
	"gitlab.com/distributed_lab/notificator-server/types"
	"github.com/Sirupsen/logrus"
)

type WorkerHerder struct {
	log *logrus.Entry
}

func NewWorkerHerder() *WorkerHerder {
	return &WorkerHerder{
		log: logrus.WithField("service", "worker_herder"),
	}
}

func (h *WorkerHerder) Init() {

}

func (h *WorkerHerder) Run() {
	workersCount := 10
	tasks := make(chan types.Request)
	results := make(chan TaskResult)
	consumed := make(chan types.Request)
	workers := make([]*Worker, workersCount)
	for i := range workers {
		workers[i] = NewWorker(tasks, results)
		go workers[i].Run()
	}

	resultsConsumer := NewTaskResultConsumer(results, consumed)
	go resultsConsumer.Run()

	taskProvider := NewTaskProvider(tasks, consumed)
	taskProvider.Run()
}
