package horizon

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/horizon-connector/internal/resources"
)

type SubmitError struct {
	statusCode int
	raw        []byte
	response   resources.SubmitResponse
}

func NewSubmitError(status int, body []byte) (SubmitError, error) {
	var submitResponse resources.SubmitResponse
	err := json.Unmarshal(body, &submitResponse)
	return SubmitError{
		statusCode: status,
		response:   submitResponse,
		raw:        body,
	}, errors.Wrap(err, "failed to unmarshal response")

}

func (e SubmitError) Error() string {
	return fmt.Sprintf("failed to submit tx: %d; tx_core: %s; op_codes: %v", e.statusCode, e.TransactionCode(), e.OperationCodes())
}

func (e *SubmitError) ResponseBody() []byte {
	return e.raw
}

func (e *SubmitError) ResponseCode() int {
	return e.statusCode
}

func (e *SubmitError) TransactionCode() string {
	return e.response.Extras.ResultCodes.Transaction
}

func (e *SubmitError) OperationCodes() []string {
	return e.response.Extras.ResultCodes.Operations
}
