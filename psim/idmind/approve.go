package idmind

import (
	"fmt"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/go/xdr"
	"context"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
)

func (s *Service) approveBothTasks(ctx context.Context, requestID uint64, requestHash string, isUSA bool) error {
	var tasksToAdd = TaskCheckIDMind
	if isUSA {
		tasksToAdd = tasksToAdd | TaskUSA
	}

	return s.approve(ctx, requestID, requestHash, tasksToAdd, TaskSubmitIDMind|TaskCheckIDMind, "{}")
}

func (s *Service) approveSubmitKYC(ctx context.Context, requestID uint64, requestHash, txID string, isUSA bool) error {
	var tasksToAdd = TaskCheckIDMind
	if isUSA {
		tasksToAdd = tasksToAdd | TaskUSA
	}

	extDetails := fmt.Sprintf(`{"%s":"%s"}`, TxIDExtDetailsKey, txID)
	return s.approve(ctx, requestID, requestHash, tasksToAdd, TaskSubmitIDMind, extDetails)
}

func (s *Service) approveCheckKYC(ctx context.Context, requestID uint64, requestHash string) error {
	// TODO In future we will probably need to add some Task at this point (e.g. for some particular admin to make some final review of final accept, or whatever)
	return s.approve(ctx, requestID, requestHash, 0, TaskCheckIDMind, "{}")
}

func (s *Service) approve(ctx context.Context, requestID uint64, requestHash string, tasksToAdd, tasksToRemove uint32, extDetails string) error {
	if extDetails == "" {
		// Just in case
		extDetails = "{}"
	}

	signedEnvelope, err := s.xdrbuilder.Transaction(s.config.Source).Op(xdrbuild.ReviewRequestOp{
		ID:     requestID,
		Hash:   requestHash,
		Action: xdr.ReviewRequestOpActionApprove,
		Details: xdrbuild.UpdateKYCDetails{
			TasksToAdd:      tasksToAdd,
			TasksToRemove:   tasksToRemove,
			ExternalDetails: extDetails,
		},
		Reason: "",
	}).Sign(s.config.Signer).Marshal()
	if err != nil {
		return errors.Wrap(err, "Failed to marshal signed Envelope")
	}

	submitResult := s.txSubmitter.Submit(ctx, signedEnvelope)
	if submitResult.Err != nil {
		return errors.Wrap(submitResult.Err, "Error submitting signed Envelope to Horizon", logan.F{
			"submit_result": submitResult,
		})
	}

	return nil
}
