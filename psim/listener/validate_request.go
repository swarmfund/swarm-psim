package listener

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
)

// ValidateHTTPRequest checks that:
//   - Request method equals to provided `requestMethod` argument,
//   - Request body is not nil,
//   - Request body is readable and
//   - (if doorman is not nil) - that request is properly signed.
//
// If any of the checks is not passed - appropriate warn log appears and appropriate response is written (true will be returned).
//
// Provide nil doorman if signature check is not needed.
//
// Note: ValidateHTTPRequest reads the body of the provided Request, so
// use returned `requestBody` bytes to parse request body.
func ValidateHTTPRequest(w http.ResponseWriter, r *http.Request, log *logan.Entry, requestMethod string, doormanChecker doorman.Doorman) (requestBody []byte, errResponseWritten bool) {
	if r.Method != requestMethod {
		log.WithField("request_method", r.Method).Warn("Received request with wrong method.")
		WriteError(w, http.StatusMethodNotAllowed, fmt.Sprintf("Only method %s is allowed.", requestMethod))
		return nil, true
	}

	if r.Body == nil {
		log.Warn("Received request with empty body.")
		WriteError(w, http.StatusBadRequest, "Empty request body.")
		return nil, true
	}

	bb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Warn("Failed to read bytes from request body Reader.")
		WriteError(w, http.StatusBadRequest, "Cannot read request body.")
		return nil, true
	}

	if doormanChecker != nil {
		var request struct {
			AccountID string `json:"account_id"`
		}

		err = json.Unmarshal(bb, &request)
		if err != nil {
			log.WithField("raw_request", string(bb)).WithError(err).Warn("Failed to preliminary unmarshal request bytes into struct(with only AccountID).")
			WriteError(w, http.StatusBadRequest, "Cannot parse JSON request.")
			return nil, true
		}

		err := doormanChecker.Check(r, doorman.SignatureOf(request.AccountID))
		if err != nil {
			log.WithError(err).Warn("Request signature is invalid.")
			WriteError(w, http.StatusUnauthorized, err.Error())
			return nil, true
		}
	}

	return bb, false
}

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

	w.Header()["Content-Type"] = append(w.Header()["Content-Type"], "application/json")
	w.WriteHeader(statusCode)

	_, err = w.Write(bb)
	if err != nil {
		return errors.Wrap(err, "Failed to write marshaled response to the ResponseWriter", logan.F{
			"marshaled_response": string(bb),
		})
	}

	return nil
}
