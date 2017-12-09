package twilio

import (
	"net/url"
	"strings"
	"fmt"
	"net/http"
)

type Connector struct {
}

func NewConnector() *Connector {
	return &Connector{
	}
}

func (c *Connector) do(payload map[string]string) (*Response, error) {
	endpoint := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", conf.SID)

	values := url.Values{}
	for k,v := range payload {
		values.Set(k, v)
	}

	rb := *strings.NewReader(values.Encode())

	client := &http.Client{}

	req, err := http.NewRequest("POST", endpoint, &rb)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(conf.SID, conf.Token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	httpResponse, err := client.Do(req)
	response := Response{
		Response: httpResponse,
	}

	return &response, err
}

func (c *Connector) SendSMS(destination, message string) (*Response, error) {

	response, err := c.do(map[string]string{
		"To": destination,
		"From": conf.FromNumber,
		"Body": message,
	})

	return response, err
}
