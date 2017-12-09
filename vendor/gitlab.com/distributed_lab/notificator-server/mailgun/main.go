package mailgun

import (
	"github.com/mailgun/mailgun-go"
)

// NewClient create new Mailgun Client from mailgun.Conf.
func NewClient() mailgun.Mailgun {
	return mailgun.NewMailgun(conf.Domain, conf.Key, conf.PublicKey)
}

// SendEmail send email throw Mailgun service.
func SendEmail(destination, subject, message string) (string, string, error) {
	mg := NewClient()
	msg := mg.NewMessage(conf.From, subject, message, destination)
	return mg.Send(msg)
}
