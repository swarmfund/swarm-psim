package server

import (
	"github.com/Sirupsen/logrus"
	"gitlab.com/distributed_lab/notificator-server/q"
	"gitlab.com/distributed_lab/notificator-server/types"
)

type TaskResultConsumer struct {
	results  <-chan TaskResult
	consumed chan<- types.Request
	log      *logrus.Entry
}

func NewTaskResultConsumer(results <-chan TaskResult, consumed chan<- types.Request, log *logrus.Entry) *TaskResultConsumer {
	return &TaskResultConsumer{
		results:  results,
		consumed: consumed,
		log:      log.WithField("service", "result_consumer"),
	}
}

func (c *TaskResultConsumer) Run() {
	for result := range c.results {
		if !result.Success {
			err := q.Request().LowerPriority(result.Request.ID)
			if err != nil {
				c.log.WithError(err).Error("failed to lower priority")
			}
		} else {
			err := q.Request().MarkCompleted(result.Request.ID)
			if err != nil {
				// actually it's pretty bad we end up here
				// task is completed but we can't update it's state
				// meaning it will be processed again eventually

				// best we can do right now is to not notify consumed channel
				// so it won't add it again during current lifespan
				c.log.WithError(err).Fatal("failed to mark request as completed")
			}
		}

		c.consumed <- result.Request
	}
}
