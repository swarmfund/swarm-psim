package mandrill

import (
	//"bullioncoin.githost.io/development/horizon/log"
	"bytes"
	"encoding/json"
	//"github.com/go-errors/errors"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	//	"log"
	"errors"

	"gitlab.com/distributed_lab/notificator-server/conf"
)

type Connector struct {
	HTTPClient  *http.Client
	API         string
	Key         string
	SenderName  string
	SenderEmail string
}

func NewConnector(cfg conf.MandrillConf) *Connector {
	return &Connector{
		HTTPClient: &http.Client{
			Timeout: time.Duration(1) * time.Minute,
		},
		API:         cfg.API,
		Key:         cfg.Key,
		SenderEmail: cfg.FromEmail,
		SenderName:  cfg.FromName,
	}
}

func (c *Connector) NotifyAccepted(receiver *Receiver, tmpl template.Template) error {
	header := "BullionCoin account accepted"
	msg := "Your account was just accepted at BullionCoin. "

	letter := Letter{header, msg}
	var buff bytes.Buffer
	err := tmpl.Execute(&buff, letter)
	if err != nil {
		//	log.WithField("err", err.Error()).Error("Error while populating template")
		return err
	}

	err = c.SendEmail(receiver, header, buff.String())
	return err
}

func (c *Connector) SendEmail(receiver *Receiver, subject, body string) error {
	message := NewMessageRequest(c.Key, NewMessage(c.SenderEmail, c.SenderName, subject, body, []*Receiver{receiver}))
	query, err := url.Parse(c.API)
	if err != nil {
		//	c.log.WithError(err).Error("Failed to parse api url")
		return err
	}
	query.Path += "/messages/send.json"
	var result json.RawMessage
	err = c.postJSON(query.String(), message, &result)
	if err != nil {
		//	c.log.WithError(err).Error("Failed to send email")
		return err
	}

	responseError := c.tryGetError(result)
	if responseError != nil {
		return responseError.ToError()
	}

	response := c.tryGetMessageResponse(result)
	if response == nil {
		return errors.New("Failed to parse response")
	}

	if response[0].RejectReason != "" {
		return errors.New("Failed to send email: " + response[0].RejectReason)
	}

	return nil
}

func (c *Connector) tryGetMessageResponse(data json.RawMessage) []MessageResponse {
	var result []MessageResponse
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil
	}

	return result
}

func (c *Connector) tryGetError(data json.RawMessage) *Error {
	var result Error
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil
	}

	return &result
}

func (c *Connector) postJSON(url string, requestBody, response interface{}) error {
	request, err := c.createPostRequest(url, requestBody)
	if err != nil {
		return err
	}

	return c.Do(request, response)
}

func (c *Connector) createPostRequest(url string, requestBody interface{}) (*http.Request, error) {
	rawRequestBody, err := c.marshalJSONBody(requestBody)
	if err != nil {
		//	c.log.WithError(err).Error("Failed to marshal request body")
		return nil, err
	}

	request, err := http.NewRequest("POST", url, rawRequestBody)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	return request, nil
}

func (c *Connector) marshalJSONBody(body interface{}) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}

	rawRequestBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(rawRequestBody), nil
}

func (c *Connector) Do(request *http.Request, response interface{}) error {
	body, err := c.do(request)
	if err != nil {
		return err
	}

	if len(body) == 0 || response == nil {
		return nil
	}

	err = json.Unmarshal(body, response)
	return err
}

func (c *Connector) do(request *http.Request) ([]byte, error) {
	request.Close = true
	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return []byte{}, err
	}

	defer func() {
		_ = resp.Body.Close()

	}()
	return ioutil.ReadAll(resp.Body)
}
