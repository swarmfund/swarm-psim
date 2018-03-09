package postage

import (
	"crypto/rand"
	"encoding/hex"

	postageapp "github.com/postageapp/postageapp-go"
)

func getUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// NewClient create new PostageApp Client from postage.Conf.
func NewClient(key string) *postageapp.Client {
	pa := new(postageapp.Client)
	pa.ApiKey = key
	return pa
}

// SendEmail send email throw PostageApp service.
func SendEmail(destination, subject, htmlMessage, from, key string) error {
	client := NewClient(key)
	msg := &postageapp.Message{
		Subject:    subject,
		From:       from,
		Html:       htmlMessage,
		Recipients: []*postageapp.Recipient{{Email: destination}},
	}

	var err error
	msg.Uid, err = getUID()
	if err != nil {
		return err
	}

	_, err = client.SendMessage(msg)
	return err
}
