package twilio

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Connector struct {
}

func NewConnector() *Connector {
	return &Connector{}
}

func (c *Connector) do(sid, token string, payload map[string]string) (*Response, error) {
	endpoint := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", sid)

	values := url.Values{}
	for k, v := range payload {
		values.Set(k, v)
	}

	rb := *strings.NewReader(values.Encode())

	client := &http.Client{}

	req, err := http.NewRequest("POST", endpoint, &rb)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(sid, token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	httpResponse, err := client.Do(req)
	response := Response{
		Response: httpResponse,
	}

	return &response, err
}

func (c *Connector) SendSMS(destination, message, number, sid, token string) (*Response, error) {

	payload := map[string]string{
		"To":   destination,
		"From": number,
		"Body": message,
	}

	response, err := c.do(sid, token, payload)

	return response, err
}
