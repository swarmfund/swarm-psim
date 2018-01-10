package scripts

import (
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"time"
	"fmt"
)

var (
	addresses = []string {
		"",
	}
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
	// TODO Get from some config
	url := ""
	authKey := ""

	// TODO Read addresses from file

	for address := range addresses {
		sendRequest(url, authKey, "importprivkey", fmt.Sprintf(`"%s", "", false`, address))
	}
}

func sendRequest(url, authKey, methodName, params string) error {
	request, err := buildRequest(url, "hardcoded_request_id", methodName, params)
	if err != nil {
		return errors.Wrap(err, "Failed to build request")
	}

	request.Header.Set("Authorization", "Basic " + authKey)

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
		return errors.Wrap(err, "Failed to unmarshal response body to JSON")
	}

	if response.Error != nil {
		return errors.Wrap(err, "Node returned non nil error")
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
