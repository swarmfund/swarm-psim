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
	req := verification.Request{}
	if ok := verification.ReadAPIRequest(s.log, w, r, &req); !ok {
		return
	}

	logger := s.log.WithFields(logan.F{
		"request": req,
	})

	envelope := xdr.TransactionEnvelope{}
	err := envelope.Scan(req.Envelope)
	if err != nil {
		logger.WithError(err).Warn("Failed to Scan TransactionEnvelope from string in request.")

		ape.RenderErr(w, r, problems.BadRequest("Cannot parse Envelope from string."))
		return
	}

	verifyErr, err := s.verifyEnvelope(envelope)
	if err != nil {
		logger.WithError(err).Error("Failed to verify TX Envelope from VerifyRequest.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}
	if verifyErr != nil {
		logger.WithField("verify_err", verifyErr).Warn("Received Request which can't pass verification - responding 403-Forbidden.")
		ape.RenderErr(w, r, problems.Forbidden(verifyErr.Error()))
		return
	}

	// TODO Render only signature in the reponse, not the whole signed Envelope, so that requester wouldn't check the returned Envelope.
	signedEnvelopeString, err := s.signAndMarshalEnvelope(envelope)
	if err != nil {
		logger.WithError(err).Error("Failed sign or marshal Envelope.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}

	ok := verification.RenderResponseEnvelope(s.log, w, r, signedEnvelopeString)
	if ok {
		logger.Info("Verified Request successfully.")
	}
}

func (s *Service) verifyEnvelope(envelope xdr.TransactionEnvelope) (verifyErr, err error) {
	if len(envelope.Tx.Operations) != 1 {
		return errors.Errorf("Must be exactly 1 Operation in the TX, found (%d).", len(envelope.Tx.Operations)), nil
	}

	opBody := envelope.Tx.Operations[0].Body

	needType := s.verifier.GetOperationType()
	if opBody.Type != needType {
		opTypeName, _ := opBody.ArmForSwitch(int32(needType))
		return errors.Errorf("Expected OperationType to be %s(%d), but got (%d).", opTypeName, needType, opBody.Type), nil
	}

	return s.verifier.VerifyOperation(envelope)
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
