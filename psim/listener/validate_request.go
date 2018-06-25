package listener

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/go/doorman"
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
