package verification

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

// ErrBadStatusFromVerifier is returned from Verify in case status code is not 2XX.
var ErrBadStatusFromVerifier = errors.New("Unsuccessful status code from Verify.")

// VerifyResponse is used to pass TX Envelope encoded to base64 between PSIMs.
type VerifyResponse struct {
	Envelope string `json:"envelope"`
}

// GetLoganFields implements logan/fields.Provider interface.
func (r VerifyResponse) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"envelope": r.Envelope,
	}
}

// Verify created verify request and sends it
// to the Verifier by provided url via http POST.
// Response is unmarshaled into VerifyResponse and
// string envelope from the response is parsed into xdr.TransactionEnvelope struct.
//
// If status code in response is not 2XX - ErrBadStatusFromVerifier is returned.
func Verify(verifierURL string, envelope string) (*xdr.TransactionEnvelope, error) {
	request := Request{
		Envelope: envelope,
	}

	return VerifyRequest(verifierURL, request)
}

// VerifyRequest marshals the provided request and sends it
// to the provided url via http POST.
// Response is unmarshaled into VerifyResponse and
// string envelope from the response is parsed into xdr.TransactionEnvelope struct.
//
// If status code in response is not 2XX - ErrBadStatusFromVerifier is returned.
//
// DEPRECATED Now this func is only used by Withdraw verifier as its implementation is not common,
// don't create new non-common implementations. Use Verify instead,
// all the data needed to be verified must be in the Envelope.
func VerifyRequest(verifierURL string, request interface{}) (*xdr.TransactionEnvelope, error) {
	rawRequestBody, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal request struct into bytes")
	}

	fields := logan.F{
		"raw_request_body": string(rawRequestBody),
	}

	bodyReader := bytes.NewReader(rawRequestBody)
	req, err := http.NewRequest("POST", verifierURL, bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create http Request to Verifier", fields)
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send created http.Request to Verifier", fields)
	}
	fields["status_code"] = resp.StatusCode

	defer func() { _ = resp.Body.Close() }()
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read the body of response from Verifier", fields)
	}
	fields["response_body"] = string(responseBytes)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.From(ErrBadStatusFromVerifier, fields)
	}

	verifyResponse := VerifyResponse{}
	err = json.Unmarshal(responseBytes, &verifyResponse)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal response body", fields)
	}

	responseEnvelope := xdr.TransactionEnvelope{}
	err = responseEnvelope.Scan(verifyResponse.Envelope)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Scan TransactionEnvelope from response from Verifier", fields)
	}

	return &responseEnvelope, nil
}
