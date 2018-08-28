package idmind

import (
	"context"

	"gitlab.com/swarmfund/psim/psim/kyc"
	"gitlab.com/tokend/regources"
)

func (s *Service) approveSubmitKYC(ctx context.Context, requestID uint64, requestHash, txID string) error {
	extDetails := map[string]string{
		TxIDExtDetailsKey: txID,
	}
	return s.requestPerformer.Approve(ctx, requestID, requestHash, 0, kyc.TaskSubmitIDMind, extDetails)
}

func (s *Service) approveCheckKYC(ctx context.Context, requestID uint64, requestHash string) error {
	// TODO In future we will probably need to add some Task at this point (e.g. for some particular admin to make some final review of final accept, or whatever)
	return s.requestPerformer.Approve(ctx, requestID, requestHash, 0, kyc.TaskCheckIDMind, nil)
}

func (s *Service) approveRequest(ctx context.Context, request regources.ReviewableRequest, externalDetails map[string]string) error {
	return s.requestPerformer.Approve(ctx, request.ID, request.Hash, 0, kyc.TaskCheckIDMind|kyc.TaskSubmitIDMind, externalDetails)
}
