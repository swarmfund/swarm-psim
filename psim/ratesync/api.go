package ratesync

import (
	"net/http"

	api "gitlab.com/tokend/psim/ape"
)

func (s *Service) API() {
	r := api.DefaultRouter()
	r.Post("/", s.VerifyHandler)
	if s.config.Pprof {
		s.log.Info("enabling debug endpoints")
		api.InjectPprof(r)
	}
	s.log.WithField("address", s.listener.Addr().String()).Info("listening")
	s.errors <- http.Serve(s.listener, r)
}

func (s *Service) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	//	defer r.Body.Close()
	//	syncRequest := SyncRequest{}
	//	err := json.NewDecoder(r.Body).Decode(&syncRequest)
	//	if err != nil {
	//		render.Render(w, r, problems.ServerErr(err))
	//		return
	//	}
	//	s.log.WithField("sync", syncRequest.Sync).Info("verification request")
	//	syncResult, ok := s.syncResults.Get(syncRequest.Sync)
	//	if !ok {
	//		s.log.WithField("sync", syncRequest.Sync).Info("don't have sync")
	//		render.Render(w, r, problems.ServerErr(errors.New("don't have sync")))
	//		return
	//	}
	//
	//	transaction := s.horizon.Transaction(&horizon.TransactionBuilder{
	//		Envelope: syncRequest.Envelope,
	//	})
	//
	//	opsVerified := 0
	//NEXT_OP:
	//	for _, op := range transaction.Operations {
	//		code := string(op.Body.ManageAssetPairOp.Base)
	//		rate := float64(op.Body.ManageAssetPairOp.PhysicalPrice)
	//		quote := op.Body.ManageAssetPairOp.Quote
	//		for _, asset := range syncResult.Assets {
	//			if string(quote) != asset.Quote {
	//				// quote assets doesn't match
	//				continue NEXT_OP
	//			}
	//			delta := math.Abs(asset.Rate*10000 - rate)
	//			if asset.Code == code && delta < 1000 {
	//				opsVerified += 1
	//				continue NEXT_OP
	//			}
	//		}
	//	}
	//
	//	if opsVerified != len(transaction.Operations) {
	//		render.Render(w, r, problems.BadRequest)
	//		return
	//	}
	//
	//	transaction.Sign(s.config.Signer)
	//
	//	err = transaction.Submit()
	//	if err != nil {
	//		entry := s.log.WithError(err)
	//		if serr, ok := err.(horizon.SubmitError); ok {
	//			entry = entry.
	//				WithField("tx code", serr.TransactionCode()).
	//				WithField("op codes", serr.OperationCodes())
	//		}
	//		entry.Error("tx failed")
	//		render.Render(w, r, problems.ServerErr(err))
	//		return
	//	}
}
