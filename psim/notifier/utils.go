package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

func sendRequest(request *http.Request, responseDest interface{}) (err error) {
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return errors.Wrap(err, "failed to make request")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			errors.Wrap(err, "Failed to close response body")
		}
	}()

	err = json.NewDecoder(resp.Body).Decode(responseDest)
	return err
}

func setJSONBody(request *http.Request, jsonBody interface{}) error {
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(jsonBody)
	if err != nil {
		return errors.Wrap(err, "failed to encode payload")
	}

	buf := b.Bytes()
	request.Body = ioutil.NopCloser(bytes.NewReader(buf))
	request.ContentLength = int64(b.Len())
	request.GetBody = func() (io.ReadCloser, error) {
		r := bytes.NewReader(buf)
		return ioutil.NopCloser(r), nil
	}
	request.Header.Set("Content-Type", "application/json")

	return nil
}
