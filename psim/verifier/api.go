package verifier

import (
	"context"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/ape/problems"
	"gitlab.com/swarmfund/psim/psim/verification"
)

// ServeAPI is blocking method.
func (s *Service) serveAPI(ctx context.Context) {
	r := ape.DefaultRouter()

	r.Post("/", s.handle)

	s.log.WithField("address", s.listener.Addr().String()).Info("Listening.")

	err := ape.ListenAndServe(ctx, s.listener, r)
	if err != nil {
		s.log.WithError(err).Error("ListenAndServe returned error.")
		return
	}
	return
}

func (s *Service) handle(w http.ResponseWriter, r *http.Request) {
	req, ok := s.verifier.ReadRequest(w, r)
	if !ok {
		return
	}

	logger := s.log.WithFields(logan.F{
		"request": req,
	})

	verifyErr, err := s.verifier.VerifyRequest(req)
	if err != nil {
		logger.WithError(err).Error("Failed to validate Request TX.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}
	if verifyErr != nil {
		logger.WithField("verify_err", verifyErr).Warn("Received Request which can't pass verification - responding 403-Forbidden.")
		ape.RenderErr(w, r, problems.Forbidden(verifyErr.Error()))
		return
	}

	signedEnvelopeString, err := s.signAndMarshalEnvelope(req.GetEnvelope())
	if err != nil {
		logger.WithError(err).Error("Failed sign or marshal Envelope.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}

	ok = verification.RenderResponseEnvelope(s.log, w, r, signedEnvelopeString)
	if ok {
		logger.Info("Verified Request successfully.")
	}
}

func (s *Service) signAndMarshalEnvelope(envelope xdr.TransactionEnvelope) (string, error) {
	fullySignedEnvelope, err := s.xdrbuilder.Sign(&envelope, s.signer)
	if err != nil {
		return "", errors.Wrap(err, "Failed to sign Envelope")
	}

	envelopeBase64, err := xdr.MarshalBase64(fullySignedEnvelope)
	if err != nil {
		return "", errors.Wrap(err, "Failed to marshal fully signed Envelope")
	}

	return envelopeBase64, nil
}
