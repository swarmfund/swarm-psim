package idmind

import (
	"net/http"

	"bytes"
	"encoding/json"
	"io/ioutil"

	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type ConnectorConfig struct {
	URL     string `fig:"url"`
	AuthKey string `fig:"auth_key"`
}

func (c ConnectorConfig) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"url": c.URL,
	}
}

type Connector struct {
	config ConnectorConfig

	client *http.Client
}

func newConnector(config ConnectorConfig) *Connector {
	return &Connector{
		config: config,

		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// Submit retrieves the data accepted by IdentityMind from KYCData,
// builds the data into the CreateAccountRequest structure
// and submits a CreateAccount request to IdentityMind.
func (c *Connector) Submit(data KYCData, email string) (*ApplicationResponse, error) {
	req, err := buildCreateAccountRequest(data, email)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create CreateAccount request to IdentityMind")
	}
	fields := logan.F{
		"create_account_request": req,
	}

	reqBB, err := json.Marshal(*req)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal CreateAccountRequest", fields)
	}

	url := c.config.URL + "/account/consumer"
	fields["url"] = url

	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(reqBB))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create HTTP request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Basic "+c.config.AuthKey)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send http POST request", fields)
	}
	fields["status_code"] = resp.StatusCode

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.From(errors.New("Unsuccessful response from IdMind"), fields)
	}

	defer func() { _ = resp.Body.Close() }()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read response body bytes", fields)
	}
	fields["response_body"] = string(respBytes)

	var appResp ApplicationResponse
	err = json.Unmarshal(respBytes, &appResp)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal response bytes into ApplicationResponse structure", fields)
	}

	return &appResp, nil
}

// TODO UploadDocument
