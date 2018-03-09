package kycairdrop

import (
	"context"
	"time"

	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/swarmfund/psim/psim/app"
)

func (s *Service) processEmails(ctx context.Context) {
	s.log.Info("Started processing emails.")

	app.RunOverIncrementalTimer(ctx, s.log, "emails_processor", func(ctx context.Context) error {
		emailsNumber := s.emails.Length()
		if emailsNumber == 0 {
			return nil
		}

		s.log.WithField("emails_number", emailsNumber).Debug("Sending emails.")

		var processedEmails []string
		s.emails.Range(ctx, func(emailAddr string) {
			logger := s.log.WithField("email_addr", emailAddr)

			err := airdrop.SendEmail(emailAddr, s.config.EmailsConfig, s.notificator)
			if err != nil {
				logger.WithError(err).Error("Failed to send email.")
				return
			}

			processedEmails = append(processedEmails, emailAddr)
			logger.Info("Notificator accepted email successfully.")
		})

		s.emails.Delete(processedEmails)
		return nil
	}, 30*time.Second, 30*time.Second)
}
