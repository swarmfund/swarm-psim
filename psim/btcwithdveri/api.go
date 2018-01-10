package btcwithdveri

import (
	"gitlab.com/swarmfund/psim/ape"
	"context"
	"encoding/json"
	"gitlab.com/swarmfund/psim/ape/problems"
	"net/http"
	"gitlab.com/swarmfund/horizon-connector"
	horizonV2 "gitlab.com/swarmfund/horizon-connector/v2"
	"fmt"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/psim/psim/btcwithdraw"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// TODO Pprof
// ServeAPI is blocking method.
func (s *Service) serveAPI(ctx context.Context) {
	r := ape.DefaultRouter()

	r.Post("/", s.handle)

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

func (s *Service) handle(w http.ResponseWriter, r *http.Request) {
	payload := btcwithdraw.ReviewRequest{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		ape.RenderErr(w, r, problems.BadRequest("Cannot parse JSON request body."))
		return
	}

	horizonTX := s.horizon.Transaction(&horizon.TransactionBuilder{
		Envelope: payload.Envelope,
	})

	if len(horizonTX.Operations) != 1 {
		ape.RenderErr(w, r, problems.BadRequest("Provided Transaction Envelope contains zero or more than one Operation."))
		return
	}

	op := horizonTX.Operations[0]
	if op.Body.Type != xdr.OperationTypeReviewRequest {
		ape.RenderErr(w, r, problems.BadRequest(fmt.Sprintf(
			"Expected Operation of type ReviewRequest(%d), but got (%d).", xdr.OperationTypeReviewRequest, op.Body.Type)))
		return
	}

	s.processReviewRequest(w, r, horizonTX)
}

func (s *Service) processReviewRequest(w http.ResponseWriter, r *http.Request, horizonTX *horizon.TransactionBuilder) {
	opBody := horizonTX.Operations[0].Body.ReviewRequestOp

	details := opBody.RequestDetails

	if details.RequestType != xdr.ReviewableRequestTypeWithdraw {
		ape.RenderErr(w, r, problems.BadRequest(fmt.Sprintf(
			"Expected Request of type Withdraw(%d), but got (%d).", xdr.ReviewableRequestTypeWithdraw, details.RequestType)))
		return
	}

	withdrawRequest, err := s.getRequest(w, r, int(opBody.RequestId))
	if err != nil {
		s.log.WithField("request_id", opBody.RequestId).WithError(err).Error("Failed to get WithdrawRequest from Horizon")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}

	if withdrawRequest.Details.RequestType != int32(xdr.ReviewableRequestTypeWithdraw) {
		// not a withdraw request
		ape.RenderErr(w, r, problems.Forbidden(fmt.Sprintf(
			"Expected Request from Horizon of type Withdraw(%d), but got (%d).", xdr.ReviewableRequestTypeWithdraw, withdrawRequest.Details.RequestType)))
		return
	}

	if withdrawRequest.State != btcwithdraw.RequestStatePending {
		ape.RenderErr(w, r, problems.Conflict(fmt.Sprintf(
			"Expected Request from Horizon with State Pending(%d), but got (%d).", btcwithdraw.RequestStatePending, withdrawRequest.State)))
		return
	}

	if withdrawRequest.Details.Withdraw.DestinationAsset != btcwithdraw.BTCAsset {
		ape.RenderErr(w, r, problems.Forbidden(fmt.Sprintf(
			"Expected Withdraw Request from Horizon with DestinationAsset %s, but got %s.", btcwithdraw.BTCAsset, withdrawRequest.Details.Withdraw.DestinationAsset)))
		return
	}

	s.processValidRequest(w, r, *withdrawRequest, horizonTX)
}

func (s *Service) getRequest(w http.ResponseWriter, r *http.Request, requestID int) (*horizonV2.Request, error){
	// TODO Stop managing requests manually, make helpers in Horizon for this instead
	req, err := s.horizon.SignedRequest("GET", fmt.Sprintf("/requests/%d", requestID), s.config.SignerKP)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create signed request (getting WithdrawalRequest)")
	}

	response, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send signed request (getting WithdrawalRequest) to Horizon")
	}

	var request horizonV2.Request
	err = json.NewDecoder(response.Body).Decode(&request)
	return &request, err
}

func (s *Service) processValidRequest(w http.ResponseWriter, r *http.Request, withdrawRequest horizonV2.Request,
		horizonTX *horizon.TransactionBuilder) {

	opBody := horizonTX.Operations[0].Body.ReviewRequestOp

	switch opBody.Action{
	case xdr.ReviewRequestOpActionApprove:
		s.processApproval(w, r, withdrawRequest, horizonTX)
		return
	case xdr.ReviewRequestOpActionPermanentReject:
		s.processReject(w, r, withdrawRequest, horizonTX)
		return
	default:
		ape.RenderErr(w, r, problems.BadRequest(fmt.Sprintf(
			"Expected Action of type Approve(%d) or PermanentReject(%d), but got (%d).",
			xdr.ReviewRequestOpActionApprove, xdr.ReviewRequestOpActionPermanentReject, opBody.Action)))
		return
	}
}
