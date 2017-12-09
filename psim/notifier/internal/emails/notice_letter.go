package emails

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"gitlab.com/tokend/go/hash"
)

const TimeLayout = "2006-01-02 15:04:05"

type NoticeTemplateType int

type NoticeLetterI interface {
	GetEmail() string
	GetHeader() string
	GetToken() string
	GetTemplateType() (string, error)
}

type NoticeLetter struct {
	ID       string
	Header   string
	Email    string
	Message  string
	Template NoticeTemplateType
}

func (letter *NoticeLetter) GetHeader() string {
	return letter.Header
}

func (letter *NoticeLetter) GetEmail() string {
	return letter.Email
}

// GetToken returns UID for letter encoded to safe base64.
func (letter *NoticeLetter) GetToken() string {
	h := hash.Hash([]byte(letter.ID))
	return base64.URLEncoding.EncodeToString(h[:])
}

func (letter *NoticeLetter) GetTemplateType() (string, error) {
	t, ok := NoticeTemplate[letter.Template]
	if !ok {
		return "", fmt.Errorf("unknown letter type: %d", letter.Template)
	}
	return t, nil
}

// ToRawMessage loads a template that matches the type of the message,
// execute it and return stringified result.
func ToRawMessage(letter NoticeLetterI) (string, error) {
	var buff bytes.Buffer
	letterType, err := letter.GetTemplateType()
	if err != nil {
		return "", err
	}

	tmpl := Templates.Lookup(letterType)
	err = tmpl.Execute(&buff, letter)
	return buff.String(), err
}
