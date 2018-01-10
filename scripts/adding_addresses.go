package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Response struct {
	ID    string `json:"id"`
	Error *Error `json:"error"`
}

type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %d message: %s", e.Code, e.Message)
}

func main() {
	log := logan.New()

	args := os.Args[1:]
	if len(args) < 2 {
		log.Panic("Need Node url(1) and auth key(2) to be passed as command line arguments.")
	}
	url := args[0]
	authKey := args[1]

	filePath := "private_keys.txt"

	privKeys, err := readPrivateKeys(filePath)
	if err != nil {
		log.WithField("file_path", filePath).WithError(err).Error("Failed to read private keys from file.")
		return
	}

	for i, privKey := range privKeys {
		if privKey == "" {
			continue
		}

		err := sendRequestToBTCNode(url, authKey, "importprivkey", fmt.Sprintf(`"%s", "", false`, privKey))
		if err != nil {
			log.WithField("i", i).WithError(err).Error("Failed to import private key.")
			return
		}

		log.WithField("i", i).Debug("Imported private key successfully.")
	}
}

func readPrivateKeys(filePath string) ([]string, error) {
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read file")
	}

	fileContent := string(dat)

	return strings.Split(fileContent, "\n"), nil
}

func sendRequestToBTCNode(url, authKey, methodName, params string) error {
	request, err := buildRequest(url, "hardcoded_request_id", methodName, params)
	if err != nil {
		return errors.Wrap(err, "Failed to build request")
	}

	request.Header.Set("Authorization", "Basic "+authKey)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(request)
	if err != nil {
		return errors.Wrap(err, "Failed to send request")
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "Failed to read response body")
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal response body to JSON", logan.F{
			"status_code":       resp.StatusCode,
			"raw_response_body": string(body),
		})
	}

	if response.Error != nil {
		return errors.Wrap(err, "Node returned non nil error", logan.F{
			"status_code": resp.StatusCode,
		})
	}

	return nil
}

func buildRequest(url, requestID, methodName, params string) (*http.Request, error) {
	bodyStr := buildRequestBody(requestID, methodName, params)
	body := bytes.NewReader([]byte(bodyStr))

	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func buildRequestBody(requestID, methodName, params string) string {
	return `{"jsonrpc": "1.0", "id":"` + requestID + `", "method": "` + methodName + `", "params": [` + params + `] }`
}
