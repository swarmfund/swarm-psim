package depositveri

import (
	"context"
	"net/http"

	"fmt"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/ape/problems"
	"gitlab.com/swarmfund/psim/psim/deposit"
	"gitlab.com/swarmfund/psim/psim/verification"
)

type request struct {
	Envelope  xdr.TransactionEnvelope
	AccountID string
}

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
	req, ok := s.readRequest(w, r)
	if !ok {
		return
	}

	logger := s.log.WithFields(logan.F{
		"account_id": req.AccountID,
	})

	envelopeFields, checkErr := s.validateDepositTX(req.Envelope)
	if checkErr != "" {
		logger.WithField("validation_err", checkErr).WithFields(envelopeFields).
			Warn("Received invalid Issuance TX - responding 403-Forbidden.")
		ape.RenderErr(w, r, problems.Forbidden(checkErr))
		return
	}

	logger = logger.WithFields(envelopeFields)

	signedEnvelopeString, err := s.signAndMarshalEnvelope(&req.Envelope)
	if err != nil {
		logger.WithError(err).Error("Failed sign or marshal Envelope.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}

	ok = verification.RenderResponseEnvelope(s.log, w, r, signedEnvelopeString)
	if ok {
		logger.Info("Verified Deposit successfully.")
	}
}

func (s *Service) readRequest(w http.ResponseWriter, r *http.Request) (*request, bool) {
	req := deposit.VerifyRequest{}
	if ok := verification.ReadAPIRequest(s.log, w, r, &req); !ok {
		return nil, false
	}

	envelope := xdr.TransactionEnvelope{}
	err := envelope.Scan(req.Envelope)
	if err != nil {
		s.log.WithField("envelope_string", req.Envelope).WithError(err).
			Warn("Failed to Scan TransactionEnvelope from string in request.")
		ape.RenderErr(w, r, problems.BadRequest("Cannot parse Envelope from string."))
		return nil, false
	}

	return &request{
		Envelope:  envelope,
		AccountID: req.AccountID,
	}, true
}

// TODO Validate that Deposit is totally correct
func (s *Service) validateDepositTX(envelope xdr.TransactionEnvelope) (envelopeFields logan.F, checkErr string) {
	fields := logan.F{}

	if len(envelope.Tx.Operations) != 1 {
		return fields, fmt.Sprintf("Must be exactly 1 Operation in the TX, found (%d).", len(envelope.Tx.Operations))
	}

	opBody := envelope.Tx.Operations[0].Body

	if opBody.Type != xdr.OperationTypeCreateIssuanceRequest {
		return fields, fmt.Sprintf("Expected OperationType to be CreateIssuanceRequest(%d), but got (%d).",
			xdr.OperationTypeCreateIssuanceRequest, opBody.Type)
	}

	op := envelope.Tx.Operations[0].Body.CreateIssuanceRequestOp

	if op == nil {
		return fields, "CreateIssuanceRequestOp is nil."
	}

	fields["reference"] = op.Reference
	fields["reference"] = op.Request

	// TODO Validate that Deposit is totally correct

	return fields, ""
}

func (s *Service) signAndMarshalEnvelope(envelope *xdr.TransactionEnvelope) (string, error) {
	fullySignedEnvelope, err := s.xdrbuilder.Sign(envelope, s.signer)
	if err != nil {
		return "", errors.Wrap(err, "Failed to sign Envelope")
	}

	envelopeBase64, err := xdr.MarshalBase64(fullySignedEnvelope)
	if err != nil {
		return "", errors.Wrap(err, "Failed to marshal fully signed Envelope")
	}

	return envelopeBase64, nil
}
