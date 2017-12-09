package notifier

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/notificator-server/client"
	"gitlab.com/swarmfund/psim/psim/notifier/internal/emails"
)

// sendEmail send email payload to notificator-server.
func (s *Service) sendEmail(letter emails.NoticeLetterI, payloadID int) error {
	if letter.GetEmail() == "" {
		s.logger.WithField("header", letter.GetHeader()).Debug("Skip sending letter cause email is empty")
		// for case when participants is not present in api db
		return nil
	}
	rawMsg, err := emails.ToRawMessage(letter)
	if err != nil {
		return err
	}
	payload := &notificator.EmailRequestPayload{
		Destination: letter.GetEmail(),
		Subject:     letter.GetHeader(),
		Message:     rawMsg,
	}

	resp, err := s.sender.Send(payloadID, letter.GetToken(), payload)
	if err != nil {
		return errors.Wrap(err, "letter is not sent")
	}

	if !resp.IsSuccess() && resp.IsPermanent() {
		return errors.From(errors.New("letter is not sent"), logan.F{
			"authenticated": resp.Authenticated(),
			"header":        letter.GetHeader(),
			"raw_response":  resp,
		})
	}

	s.logger.WithFields(logan.F{
		"success": resp.IsSuccess(),
		"header":  letter.GetHeader(),
	}).Debug("Letter sent.")

	return err
}
