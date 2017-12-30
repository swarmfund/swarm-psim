package btcwithdveri

import (
	"gitlab.com/swarmfund/psim/ape"
	"context"
	"encoding/json"
	"gitlab.com/swarmfund/psim/ape/problems"
	"net/http"
	"gitlab.com/swarmfund/horizon-connector"
)

func (s *Service) serveAPI(ctx context.Context) {
	r := ape.DefaultRouter()

	r.Post("/approve", s.approveHandler)
	r.Post("/reject", s.rejectHandler)

	// TODO
	//if s.config.Pprof {
	//	s.log.Info("enabling debugging endpoints")
	//	ape.InjectPprof(r)
	//}

	s.log.WithField("address", s.listener.Addr().String()).Info("listening")

	err := ape.ListenAndServe(ctx, s.listener, r)
	if err != nil {
		s.log.WithError(err).Error("ListenAndServe returned error.")
		return
	}
	return
}

type ApproveRequest struct {
	Envelope string `json:"envelope"`

	RequestID uint64 `json:"request_id"`
	RequestHash uint64 `json:"request_hash"`
}

// TODO Call verifyApprove() method to make the job done.
// TODO split me to several methods
func (s *Service) approveHandler(w http.ResponseWriter, r *http.Request) {
	payload := ApproveRequest{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		ape.RenderErr(w, r, problems.BadRequest("Cannot parse JSON request."))
		return
	}

	horizonTX := s.horizon.Transaction(&horizon.TransactionBuilder{
		Envelope: payload.Envelope,
	})

	if len(horizonTX.Operations) != 1 {
		ape.RenderErr(w, r, problems.BadRequest("Provided Transaction Envelope contains more than one Operation."))
		return
	}

	// TODO Validate Envelope
	//op := horizonTX.Operations[0]
	//if op.Body.Type != xdr.OperationTypeManageCoinsEmissionRequest {
	//	ape.RenderErr(w, r, problems.BadRequest(fmt.Sprintf(
	//		"Expected Operation of type ManageCoinEmissionRequest(%d), but got (%d).",
	//		xdr.OperationTypeManageCoinsEmissionRequest, op.Body.Type)))
	//	return
	//}
	//
	//opBody := horizonTX.Operations[0].Body.ManageCoinsEmissionRequestOp
	//
	//if opBody.Asset != bitcoin.Asset {
	//	ape.RenderErr(w, r, problems.BadRequest(fmt.Sprintf("Expected asset to be '%s', but got '%s'.",
	//		bitcoin.Asset, opBody.Asset)))
	//	return
	//}
	//
	//if opBody.Action != xdr.ManageCoinsEmissionRequestActionManageCoinsEmissionRequestCreate {
	//	ape.RenderErr(w, r, problems.BadRequest("Expected Action to be CER create."))
	//	return
	//}
	//
	//reference := bitcoin.BuildCoinEmissionRequestReference(payload.TXHash, payload.OutIndex)
	//if string(opBody.Reference) != reference {
	//	ape.RenderErr(w, r, problems.Conflict(fmt.Sprintf("Expected reference to be '%s', but got '%s'.",
	//		reference, string(opBody.Reference))))
	//	return
	//}
}

// TODO
func (s *Service) rejectHandler(w http.ResponseWriter, r *http.Request) {

}
