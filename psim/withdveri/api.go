package withdveri

import (
	"context"
	"encoding/json"
	"net/http"

	"fmt"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/ape/problems"
	"gitlab.com/swarmfund/psim/psim/withdraw"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/verification"
)

// TODO Pprof
// ServeAPI is blocking method.
func (s *Service) serveAPI(ctx context.Context) {
	r := ape.DefaultRouter()

	r.Post(withdraw.VerifyPreliminaryApproveURLSuffix, s.preliminaryApproveHandler)
	r.Post(withdraw.VerifyApproveURLSuffix, s.approveHandler)
	r.Post(withdraw.VerifyRejectURLSuffix, s.rejectHandler)

	// TODO
	//if s.config.Pprof {
	//	s.log.Info("enabling debugging endpoints")
	//	ape.InjectPprof(r)
	//}

	s.log.WithField("address", s.listener.Addr().String()).Info("Listening.")

	err := ape.ListenAndServe(ctx, s.listener, r)
	if err != nil {
		s.log.WithError(err).Error("ListenAndServe returned error.")
		return
	}
	return
}

func (s *Service) obtainAndCheckRequest(requestID uint64, requestHash string, neededRequestType int32) (request *horizon.Request, checkErr string, err error) {
	request, err = withdraw.ObtainRequest(s.horizon.Client(), requestID)
	if err != nil {
		return nil, "", errors.Wrap(err, "Failed to Obtain WithdrawRequest from Horizon")
	}

	if request.Hash != requestHash {
		return nil, fmt.Sprintf("The RequestHash from Horizon (%s) does not match the one provided (%s).", request.Hash, requestHash), nil
	}
	proveErr := withdraw.ProvePendingRequest(*request, s.offchainHelper.GetAsset(), neededRequestType)
	if proveErr != "" {
		return nil, fmt.Sprintf("Not a pending WithdrawRequest: %s", proveErr), nil
	}

	return request, "", nil
}
