package server

import (
	"time"
	"gitlab.com/distributed_lab/notificator-server/q"
	"gitlab.com/distributed_lab/notificator-server/types"
	"gitlab.com/distributed_lab/notificator-server/log"
)

type TaskProvider struct {
	tasks    chan<- types.Request
	consumed <-chan types.Request

	// TODO wrap access with locks
	processing map[int64]bool
}

func (p *TaskProvider) Run() {
	entry := log.WithField("service", "task_provider")
	tick := time.NewTicker(5 * time.Second)
	go func() {
		for request := range p.consumed {
			entry.WithField("request", request.ID).Info("consumed")
			delete(p.processing, request.ID)
		}
	}()
	for {
		requests, err := q.Request().GetHead()
		if err != nil {
			entry.WithError(err).Error("failed to get requests")
			continue
		}

		entry.WithField("count", len(requests)).Info("adding requests")
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

func NewTaskProvider(tasks chan<- types.Request, consumed <-chan types.Request) *TaskProvider {
	return &TaskProvider{
		tasks:      tasks,
		consumed:   consumed,
		processing: map[int64]bool{},
	}
}
