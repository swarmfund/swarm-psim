package idmind

import (
	"net/http"

	"bytes"
	"encoding/json"
	"io/ioutil"

	"fmt"
	"io"
	"mime/multipart"

	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"strings"
)

var ErrAppNotFound = errors.New("IDMind failed to find Application with the provided TXId.")

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

		client: &http.Client{Timeout: 2 * time.Minute},
	}
}

// Submit retrieves the data accepted by IdentityMind from KYCData,
// builds the data into the CreateAccountRequest structure
// and submits a CreateAccount request to IdentityMind.
func (c *Connector) Submit(req CreateAccountRequest) (*ApplicationResponse, error) {
	url := c.config.URL + "/account/consumer"
	fields := logan.F{
		"url": url,
	}

	reqBB, err := json.Marshal(req)
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

	defer func() { _ = resp.Body.Close() }()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read response body bytes", fields)
	}
	fields["response_body"] = string(respBytes)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.From(errors.New("Unsuccessful response from IdMind"), fields)
	}

	var appResp ApplicationResponse
	err = json.Unmarshal(respBytes, &appResp)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal response bytes into ApplicationResponse structure", fields)
	}

	return &appResp, nil
}

func (c *Connector) UploadDocument(txID, description string, fileName string, fileReader io.Reader) error {
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

	// This code is commented intentionally. IDMind says this field is required, but Sandbox works fine without it :/
	//err = bufferWriter.WriteField("appID", "424284")
	//if err != nil {
	//	return errors.Wrap(err, "Failed to add appID field to request")
	//}

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

	req.Header.Set("Authorization", "Basic "+c.config.AuthKey)
	req.Header.Set("Content-Type", bufferWriter.FormDataContentType())

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

func (c *Connector) CheckState(txID string) (*CheckApplicationResponse, error) {
	url := fmt.Sprintf("%s/account/consumer/%s", c.config.URL, txID)
	fields := logan.F{
		"url": url,
	}

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create HTTP request")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Basic "+c.config.AuthKey)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send http GET request", fields)
	}
	fields["status_code"] = resp.StatusCode

	defer func() { _ = resp.Body.Close() }()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read response body bytes", fields)
	}
	respS := string(respBytes)
	fields["raw_response_body"] = respS

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// The way of detecting this particular error from IDMind is quite dirty, but they don't have any error codes to parse their errors more strictly.
		if resp.StatusCode == http.StatusBadRequest && strings.Contains(respS, fmt.Sprintf(`{"error_message":"Failed to find application with id: %s"}`, txID)) {
			return nil, ErrAppNotFound
		}

		return nil, errors.From(errors.New("Unsuccessful response from IdMind"), fields)
	}

	var checkResp CheckApplicationResponse
	err = json.Unmarshal(respBytes, &checkResp)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal response bytes into CheckApplicationResponse structure", fields)
	}

	return &checkResp, nil
}
