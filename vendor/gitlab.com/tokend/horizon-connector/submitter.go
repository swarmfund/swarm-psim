package horizon

import (
	"bytes"
	"context"
	"encoding/json"

	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon-connector/internal/resources"
	"gitlab.com/tokend/horizon-connector/internal/responses"
)

var (
	ErrSubmitTimeout              = errors.New("submit timed out")
	ErrSubmitInternal             = errors.New("internal submit error")
	ErrSubmitRejected             = errors.New("transaction rejected")
	ErrSubmitMalformed            = errors.New("transaction malformed")
	ErrSubmitUnexpectedStatusCode = errors.New("Unexpected unsuccessful status code.")
)

type Submitter struct {
	client *Client
}

type SubmitResult struct {
	Err         error
	RawResponse []byte
	TXCode      string
	OpCodes     []string
}

func (r SubmitResult) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"err":      r.Err.Error(),
		"raw":      string(r.RawResponse),
		"tx_code":  r.TXCode,
		"op_codes": r.OpCodes,
	}
}

// DEPRECATED
// Use SubmitE method instead
func (s *Submitter) Submit(ctx context.Context, envelope string) SubmitResult {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(&resources.TransactionSubmit{
		Transaction: envelope,
	})
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal request"))
	}
	response, err := s.client.Post("/transactions", &buf)
	if err == nil {
		// successful submission
		return SubmitResult{
			RawResponse: response,
		}
	}
	cerr := errors.Cause(err).(Error)
	result := SubmitResult{
		RawResponse: cerr.Body(),
	}
	// go through known response codes and try to build meaningful result
	switch cerr.Status() {
	case http.StatusGatewayTimeout: // timeout
		result.Err = ErrSubmitTimeout
	case http.StatusBadRequest: // rejected or malformed
		// check which error it was exactly, might be useful for consumer
		var response responses.TransactionBadRequest
		if err := json.Unmarshal(result.RawResponse, &response); err != nil {
			panic(errors.Wrap(err, "failed to unmarshal horizon response"))
		}
		switch response.Type {
		case "transaction_malformed":
			result.Err = ErrSubmitMalformed
		case "transaction_failed":
			result.Err = ErrSubmitRejected
			result.TXCode = response.Extras.ResultCodes.Transaction
			result.OpCodes = response.Extras.ResultCodes.Operations
		default:
			panic("unknown reject type")
		}
	case http.StatusInternalServerError: // internal error
		result.Err = ErrSubmitInternal
	default:
		// TODO poke someone who touched horizon
		panic("unexpected submission result")
	}
	return result
}

// SubmitE submits txEnvelope checking common unsuccessful status codes and forming
// appropriate errors for status codes.
// For transaction_failed response in BadRequest(400) response, SubmitE parses
// TXCode and OpCodes, returning them with the SubmitFailedError.
//
// Consumer of SubmitE can view the response TxCode and OpCodes, if transaction_failed reject response
func (s *Submitter) SubmitE(txEnvelope string) error {
	req := resources.TransactionSubmit{
		Transaction: txEnvelope,
	}

	statusCode, respBB, err := s.client.PostJSON("/transactions", req)
	if err != nil {
		return errors.Wrap(err, "Failed to send POST request via Client")
	}
	fields := logan.F{
		"status_code":  statusCode,
		"raw_response": string(respBB),
	}

	if isStatusCodeSuccessful(statusCode) {
		return nil
	}

	// go through known unsuccessful response status codes and try to build meaningful result
	switch statusCode {
	case http.StatusGatewayTimeout:
		return ErrSubmitTimeout
	case http.StatusBadRequest: // rejected or malformed
		// check which error it was exactly, might be useful for consumer
		var response responses.TransactionBadRequest
		if err := json.Unmarshal(respBB, &response); err != nil {
			return errors.Wrap(err, "Failed to unmarshal BadRequest Horizon response bytes into TransactionBadRequest struct", fields)
		}

		switch response.Type {
		case "transaction_malformed":
			return errors.From(ErrSubmitMalformed, fields)
		case "transaction_failed":
			return NewSubmitFailedError("Transaction failed.", statusCode, response.Extras.ResultCodes.Transaction, response.Extras.ResultCodes.Operations)
		default:
			return errors.From(errors.New("Unknown reject type."), fields.Merge(logan.F{
				"reject_type": response.Type,
			}))
		}
	case http.StatusInternalServerError: // internal error
		return errors.From(ErrSubmitInternal, fields)
	default:
		// Normally must never happen. Looks like somebody changed Horizon.
		return errors.From(ErrSubmitUnexpectedStatusCode, fields)
	}

	return nil
}
