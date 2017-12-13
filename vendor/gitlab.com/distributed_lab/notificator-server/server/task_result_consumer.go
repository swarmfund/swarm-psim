package server

import (
	"gitlab.com/distributed_lab/notificator-server/q"
	"gitlab.com/distributed_lab/notificator-server/types"
	"gitlab.com/distributed_lab/notificator-server/log"
)

type TaskResultConsumer struct {
	results  <-chan TaskResult
	consumed chan<- types.Request
}

func NewTaskResultConsumer(results <-chan TaskResult, consumed chan<- types.Request) *TaskResultConsumer {
	return &TaskResultConsumer{
		results:  results,
		consumed: consumed,
	}
}

func (c *TaskResultConsumer) Run() {
	entry := log.WithField("service", "result_consumer")
	for result := range c.results {
		if !result.Success {
			err := q.Request().LowerPriority(result.Request.ID)
			if err != nil {
				entry.WithError(err).Error("failed to lower priority")
			}
		} else {
			err := q.Request().MarkCompleted(result.Request.ID)
			if err != nil {
				// actually it's pretty bad we end up here
				// task is completed but we can't update it's state
				// meaning it will be processed again eventually

				// best we can do right now is to not notify consumed channel
				// so it won't add it again during current lifespan
				entry.WithError(err).Fatal("failed to mark request as completed")
			}
		}

		c.consumed <- result.Request
	}
}
