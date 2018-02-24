package airdrop

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"
	notificator "gitlab.com/distributed_lab/notificator-server/client"
)

func (s *Service) processEmails(ctx context.Context) {
	ticker := time.Tick(30 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker:
			if s.emails.Length() == 0 {
				break
			}

			var processedEmails []string
			s.emails.Range(ctx, func(email string) {
				err := s.sendEmail(email)
				if err != nil {
					s.log.WithField("email", email).WithError(err).Error("Failed to send email.")
					return
				}

				processedEmails = append(processedEmails, email)
			})

			s.emails.Delete(processedEmails)
		}
	}
}

func (s *Service) sendEmail(email string) error {
	msg := &notificator.EmailRequestPayload{
		Destination: email,
		// TODO
		Subject: "My awesome subject",
		// TODO
		Message: "My awesome message.",
	}

	resp, err := s.notificator.Send(s.config.EmailRequestType, email, msg)
	if err != nil {
		return errors.Wrap(err, "Failed to send email via Notificator")
	}

	if !resp.IsSuccess() {
		return errors.New("Unsuccessful email sending.")
	}

	return nil
}
