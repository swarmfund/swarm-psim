package mailgun

import (
	"github.com/mailgun/mailgun-go"
)

// NewClient create new Mailgun Client from mailgun.Conf.
func NewClient(domain, key, publicKey string) mailgun.Mailgun {
	return mailgun.NewMailgun(domain, key, publicKey)
}

// SendEmail send email throw Mailgun service.
func SendEmail(destination, subject, message, from, domain, key, publicKey string) (string, string, error) {
	mg := NewClient(domain, key, publicKey)
	msg := mg.NewMessage(from, subject, "", destination)
	msg.SetHtml(message)
	return mg.Send(msg)
}
