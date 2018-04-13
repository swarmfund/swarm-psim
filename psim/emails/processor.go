package emails

import (
	"net/http"

	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/notificator-server/client"
	"gitlab.com/distributed_lab/running"
)

type NotificatorConnector interface {
	Send(requestType int, token string, payload notificator.Payload) (*notificator.Response, error)
}

// Processor makes retries of sending emails using Notificator for you.
// Processor can only work with single subject and email message,
// which are provided in config to the constructor.
type Processor struct {
	log         *logan.Entry
	config      Config
	notificator NotificatorConnector

	tasks TaskSyncSet
}

// NewProcessor is just a constructor for Processor
func NewProcessor(
	log *logan.Entry,
	config Config,
	notificator NotificatorConnector) *Processor {

	return &Processor{
		log:         log.WithField("helper-runner", "emails_processor"),
		config:      config,
		notificator: notificator,

		tasks: newSyncSet(),
	}
}

// Run is locking function, returns only when ctx cancels.
func (p *Processor) Run(ctx context.Context) {
	p.log.WithField("", p.config).Info("Started emails processor.")

	running.WithBackOff(ctx, p.log, "emails_processor", func(ctx context.Context) error {
		tasksNumber := p.tasks.length()
		if tasksNumber == 0 {
			p.log.Debugf("No emails to send - waiting for next wake up (%s).", p.config.SendPeriod.String())
			return nil
		}

		p.log.WithField("tasks_number", tasksNumber).Debug("Sending emails.")

		var processedTasks []Task
		p.tasks.rangeThrough(ctx, func(task Task) {
			logger := p.log.WithField("task", task)

			emailWasSent, err := p.sendEmail(task.toPayload())
			if err != nil {
				logger.WithError(err).Error("Failed to send email.")
				return
			}

			processedTasks = append(processedTasks, task)

			if emailWasSent {
				logger.Info("Notificator accepted email task successfully.")
			} else {
				logger.Debug("Email has been already sent earlier - skipping.")
			}
		})

		p.tasks.delete(processedTasks)
		return nil
	}, p.config.SendPeriod, 30*time.Second, time.Hour)
}

func (p *Processor) AddEmailAddresses(ctx context.Context, subject, message string, emailAddresses []string) {
	for _, addr := range emailAddresses {
		p.tasks.put(ctx, Task{
			Destination: addr,
			Subject:     subject,
			Message:     message,
		})
	}
}

func (p *Processor) AddTask(ctx context.Context, emailAddress, subject, message string) {
	p.AddEmailAddresses(ctx, subject, message, []string{emailAddress})
}

func (p *Processor) sendEmail(payload notificator.EmailRequestPayload) (bool, error) {
	uniqueToken := payload.Destination + p.config.UniquenessTokenSuffix

	resp, err := p.notificator.Send(p.config.RequestType, uniqueToken, payload)
	if err != nil {
		return false, errors.Wrap(err, "Failed to send email via Notificator")
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		// The emails has already been sent earlier.
		return false, nil
	}

	if !resp.IsSuccess() {
		return false, errors.From(errors.New("Unsuccessful response for email sending request."), logan.F{
			"notificator_response": resp,
			"status_code":          resp.StatusCode,
		})
	}

	return true, nil
}
