package idmind

import (
	"net/http"

	"bytes"
	"encoding/json"
	"io/ioutil"

	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"fmt"
	"mime/multipart"
	"io"
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
	url := c.config.URL + "/account/consumer"
	fields := logan.F{
		"url": url,
	}

	req, err := buildCreateAccountRequest(data, email)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create CreateAccount request to IdentityMind")
	}
	fields["create_account_request"] = req

	reqBB, err := json.Marshal(*req)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal CreateAccountRequest", fields)
	}

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

func (c *Connector) UploadDocument(appID, txID, description string, fileName string, fileReader io.Reader) error {
	url := fmt.Sprintf("%s/account/consumer/%s/files", c.config.URL, txID)

	var buffer bytes.Buffer
	bufferWriter := multipart.NewWriter(&buffer)

	// Write file
	fileWriter, err := bufferWriter.CreateFormFile("file", fileName)
	if err != nil {
		return errors.Wrap(err, "Failed to create bufferWriter from file")
	}

	_, err = io.Copy(fileWriter, fileReader)
	if err != nil {
		return errors.Wrap(err, "Failed to copy file from Reader to Writer")
	}

	// Write simple fields
	err = bufferWriter.WriteField("appID", appID)
	if err != nil {
		return errors.Wrap(err, "Failed to add appID field to request")
	}
	if description != "" {
		err = bufferWriter.WriteField("description", description)
		if err != nil {
			return errors.Wrap(err, "Failed to add description field to request")
		}
	}

	err = bufferWriter.Close()
	if err != nil {
		return errors.Wrap(err, "Failed to close request bufferWriter")
	}

	req, err := http.NewRequest("POST", url, &buffer)
	if err != nil {
		return errors.Wrap(err, "Failed to create http request")
	}
	//bufferWriter.FormDataContentType()

	req.Header.Set("Authorization", "Basic "+c.config.AuthKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "Failed to send http POST request")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.From(errors.New("Unsuccessful response from IdMind"), logan.F{
			"status_code": resp.StatusCode,
		})
	}

	return nil
}
