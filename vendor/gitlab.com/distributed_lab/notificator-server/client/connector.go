package notificator

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

var (
	ErrInternal = errors.New("internal error")
)

type Connector struct {
	client   *http.Client
	pair     Pair
	endpoint url.URL
}


func NewConnector(pair Pair, endpoint url.URL) *Connector {
	return &Connector{
		client:   http.DefaultClient,
		pair:     pair,
		endpoint: endpoint,
	}
}

func (c *Connector) Send(requestType int, token string, payload Payload) (*Response, error) {
	data := map[string]interface{}{
		"type":    requestType,
		"token":   token,
		"payload": payload,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return nil, ErrInternal
	}


	request, err := http.NewRequest("POST", c.endpoint.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, ErrInternal
	}

	signature := c.pair.Signature(body)
	request.Header.Set("authorization", c.pair.Public)
	request.Header.Set("x-signature", signature)

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	apiResponse := new(apiResponse)
	_ = json.NewDecoder(response.Body).Decode(&apiResponse)

	return &Response{
		statusCode:  response.StatusCode,
		apiResponse: apiResponse,
	}, nil
}
