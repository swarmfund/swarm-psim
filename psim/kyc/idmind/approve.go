package idmind

import (
	"context"
	"gitlab.com/swarmfund/psim/psim/kyc"
)

func (s *Service) approveBothTasks(ctx context.Context, requestID uint64, requestHash string, isUSA bool) error {
	var tasksToAdd uint32
	if isUSA {
		tasksToAdd = tasksToAdd | kyc.TaskUSA
	}

	return s.requestPerformer.Approve(ctx, requestID, requestHash, tasksToAdd, kyc.TaskSubmitIDMind|kyc.TaskCheckIDMind|kyc.TaskNonLatinDoc, nil)
}

func (s *Service) approveSubmitKYC(ctx context.Context, requestID uint64, requestHash, txID string) error {
	extDetails := map[string]string {
		TxIDExtDetailsKey: txID,
	}
	return s.requestPerformer.Approve(ctx, requestID, requestHash, 0, kyc.TaskSubmitIDMind, extDetails)
}

func (s *Service) approveCheckKYC(ctx context.Context, requestID uint64, requestHash string) error {
	// TODO In future we will probably need to add some Task at this point (e.g. for some particular admin to make some final review of final accept, or whatever)
	return s.requestPerformer.Approve(ctx, requestID, requestHash, 0, kyc.TaskCheckIDMind, nil)
}
