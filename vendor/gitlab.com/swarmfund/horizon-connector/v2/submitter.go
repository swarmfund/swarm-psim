package horizon

import (
	"bytes"
	"context"
	"encoding/json"

	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/resources"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/responses"
)

var (
	ErrSubmitTimeout   = errors.New("submit timed out")
	ErrSubmitInternal  = errors.New("internal submit error")
	ErrSubmitRejected  = errors.New("transaction rejected")
	ErrSubmitMalformed = errors.New("transaction malformed")
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
	return map[string]interface{} {
		"err":      r.Err.Error(),
		"raw":      string(r.RawResponse),
		"tx_code":  r.TXCode,
		"op_codes": r.OpCodes,
	}
}

// TODO Return error
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
