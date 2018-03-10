package server

import (
	"time"

	"github.com/Sirupsen/logrus"
	"gitlab.com/distributed_lab/notificator-server/q"
	"gitlab.com/distributed_lab/notificator-server/types"
)

type TaskProvider struct {
	tasks    chan<- types.Request
	consumed <-chan types.Request
	log      *logrus.Entry
	// TODO wrap access with locks
	processing map[int64]bool
}

func (p *TaskProvider) Run() {
	tick := time.NewTicker(5 * time.Second)
	go func() {
		for request := range p.consumed {
			p.log.WithField("request", request.ID).Info("consumed")
			delete(p.processing, request.ID)
		}
	}()
	for {
		requests, err := q.Request().GetHead()
		if err != nil {
			p.log.WithError(err).Error("failed to get requests")
			continue
		}

		p.log.WithField("count", len(requests)).Info("adding requests")
		for _, request := range requests {
			if _, ok := p.processing[request.ID]; ok {
				continue
			}
			p.tasks <- request
			p.processing[request.ID] = true
		}

		<-tick.C
	}
}

func NewTaskProvider(tasks chan<- types.Request, consumed <-chan types.Request, log *logrus.Entry) *TaskProvider {
	return &TaskProvider{
		tasks:      tasks,
		consumed:   consumed,
		log:        log.WithField("service", "task_provider"),
		processing: map[int64]bool{},
	}
}
