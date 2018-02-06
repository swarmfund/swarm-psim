package verification

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/tokend/keypair"
)

// ErrBadStatusFromVerifier is returned from SendRequestToVerifier in case status code is not 2XX.
var ErrBadStatusFromVerifier = errors.New("Unsuccessful status code from Verify.")

// EnvelopeResponse is used to pass TX Envelope encoded to base64 between PSIMs.
type EnvelopeResponse struct {
	Envelope string `json:"envelope"`
}

// GetLoganFields implements logan/fields.Provider interface.
func (r EnvelopeResponse) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"envelope": r.Envelope,
	}
}

// SendRequestToVerifier marshals the provided request and sends it
// to the provided url via http POST.
// Response is unmarshaled into EnvelopeResponse and
// string envelope from the response is parsed into xdr.TransactionEnvelope struct.
//
// If status code in response is not 2XX - ErrBadStatusFromVerifier is returned.
func SendRequestToVerifier(url string, request interface{}) (*xdr.TransactionEnvelope, error) {
	rawRequestBody, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal request struct into bytes")
	}

	fields := logan.F{
		"verify_url":       url,
		"raw_request_body": string(rawRequestBody),
	}

	bodyReader := bytes.NewReader(rawRequestBody)
	req, err := http.NewRequest("POST", url, bodyReader)
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

	envelopeResponse := EnvelopeResponse{}
	err = json.Unmarshal(responseBytes, &envelopeResponse)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal response body", fields)
	}

	envelope := xdr.TransactionEnvelope{}
	err = envelope.Scan(envelopeResponse.Envelope)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Scan TransactionEnvelope from response from Verifier", fields)
	}

	return &envelope, nil
}
