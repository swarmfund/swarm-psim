package btcwithdveri

import (
	"context"
	"gitlab.com/swarmfund/horizon-connector"
	horizonV2 "gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/btcwithdraw"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/ape/problems"
	"net/http"
	"gitlab.com/distributed_lab/logan/v3"
	"fmt"
	"github.com/go-errors/errors"
)

func (s *Service) processReject(w http.ResponseWriter, r *http.Request, withdrawRequest *horizonV2.Request, horizonTX *horizon.TransactionBuilder) {
	ctx := r.Context()
	opBody := horizonTX.Operations[0].Body.ReviewRequestOp
	fields := logan.F{
		"request_id": opBody.RequestId,
		"request_hash": opBody.RequestHash,
		"request_action_i": opBody.Action,
		"request_action": opBody.Action.String(),
		"reject_reason": opBody.Reason,
	}

	validationErr := s.validateReject(withdrawRequest, btcwithdraw.RejectReason(opBody.Reason))
	if validationErr != "" {
		ape.RenderErr(w, r, problems.Forbidden(validationErr))
		return
	}

	err := s.processValidReject(ctx, horizonTX)
	if err != nil {
		s.log.WithFields(fields).WithError(err).Error("Failed to process valid PermanentReject Request.")
		ape.RenderErr(w, r, problems.ServerError(err))
	}

	s.log.WithFields(fields).Info("Verified PermanentReject for WithdrawalRequest successfully.")
}

// TODO
func (s *Service) validateReject(withdrawRequest *horizonV2.Request, reason btcwithdraw.RejectReason) string {
	switch reason {
	case btcwithdraw.RejectReasonInvalidAddress:
		// TODO
	case btcwithdraw.RejectReasonTooLittleAmount:
		// TODO
	}

	return fmt.Sprintf("Unsupported RejectReason '%s'.", string(reason))
}

// TODO
func (s *Service) processValidReject(ctx context.Context, tx *horizon.TransactionBuilder) error {
	// TODO
	return errors.New("Not implemented.")
}
