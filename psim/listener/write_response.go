package listener

import (
	"net/http"
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
)

// WriteError takes the errorMessage string, puts it into `error` field of a JSON,
// marshals this JSON and writes it into the response.
// Possible errors of marshal or response writer write can be returned.
func WriteError(w http.ResponseWriter, statusCode int, errorMessage string) error {
	resp := struct {
		Error string `json:"error"`
	}{
		Error: errorMessage,
	}

	bb, err := json.Marshal(resp)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal response to bytes")
	}

	// Wrapping of nil error returns nil.
	return errors.Wrap(WriteResponse(w, statusCode, bb), "Failed to write response")
}

// WriteResponse won't try to write response body into the ResponseWriter, if responseBody argument is nil.
func WriteResponse(w http.ResponseWriter, statusCode int, responseBody []byte) error {
	w.Header()["Content-Type"] = append(w.Header()["Content-Type"], "application/json")
	w.WriteHeader(statusCode)

	if responseBody != nil {
		_, err := w.Write(responseBody)
		if err != nil {
			return errors.Wrap(err, "Failed to write marshaled response into the ResponseWriter", logan.F{
				"marshaled_response": string(responseBody),
			})
		}
	}

	return nil
}
