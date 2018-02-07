package depositveri

import (
	"context"
	"net/http"

	"fmt"

	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/ape/problems"
	"gitlab.com/swarmfund/psim/psim/deposit"
	"gitlab.com/swarmfund/psim/psim/verification"
)

var errNoExtAccount = errors.New("External system Account was not found.")

type request struct {
	AccountID string
	Envelope  xdr.TransactionEnvelope
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
		"account_id":           req.AccountID,
		"envelope_ext_details": getExternalDetails(req.Envelope),
	})

	checkErr, err := s.validateDeposit(*req)
	if err != nil {
		logger.WithError(err).Error("Failed to validate Deposit TX.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}
	if checkErr != "" {
		logger.WithField("validation_err", checkErr).Warn("Received invalid Issuance TX - responding 403-Forbidden.")
		ape.RenderErr(w, r, problems.Forbidden(checkErr))
		return
	}

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

// GetExternalDetails is safe - it returns empty string if any shit happens(normally it doesn't happen).
func getExternalDetails(envelope xdr.TransactionEnvelope) string {
	if len(envelope.Tx.Operations) == 0 {
		return ""
	}

	opBody := envelope.Tx.Operations[0].Body

	if opBody.CreateIssuanceRequestOp == nil {
		return ""
	}

	return string(opBody.CreateIssuanceRequestOp.Request.ExternalDetails)
}

func (s *Service) validateDeposit(req request) (checkErr string, err error) {
	if len(req.Envelope.Tx.Operations) != 1 {
		return fmt.Sprintf("Must be exactly 1 Operation in the TX, found (%d).", len(req.Envelope.Tx.Operations)), nil
	}

	opBody := req.Envelope.Tx.Operations[0].Body

	if opBody.Type != xdr.OperationTypeCreateIssuanceRequest {
		return fmt.Sprintf("Expected OperationType to be CreateIssuanceRequest(%d), but got (%d).",
			xdr.OperationTypeCreateIssuanceRequest, opBody.Type), nil
	}

	op := req.Envelope.Tx.Operations[0].Body.CreateIssuanceRequestOp

	if op == nil {
		return "CreateIssuanceRequestOp is nil.", nil
	}

	checkErr, err = s.validateIssuanceOp(*op, req.AccountID)
	if err != nil {
		return "", errors.Wrap(err, "Failed to validate Issuance Op")
	}

	return checkErr, nil
}

// TODO Try make me smaller
func (s *Service) validateIssuanceOp(op xdr.CreateIssuanceRequestOp, accountAddress string) (checkErr string, err error) {
	req := op.Request

	if string(req.Asset) != s.offchainHelper.GetAsset() {
		return fmt.Sprintf("Invalid asset - expected (%s), got (%s).", s.offchainHelper.GetAsset(), req.Asset), nil
	}

	balanceID, err := deposit.GetBalanceID(s.horizon, accountAddress, s.offchainHelper.GetAsset())
	if err != nil {
		return "", errors.Wrap(err, "Failed to get BalanceID by AccountAddress and Asset from Horizon")
	}
	if req.Receiver.AsString() != balanceID {
		return fmt.Sprintf("Invalid BalanceID - expected (%s), got (%s).", balanceID, req.Receiver.AsString()), nil
	}

	// TODO move "bitcoin" to some config
	offchainAddress, err := s.getOffchainAddress(accountAddress, "bitcoin")
	if err != nil {
		return "", errors.Wrap(err, "Failed to get Offchain Address of the Account")
	}

	extDetails := deposit.ExternalDetails{}
	err = json.Unmarshal([]byte(string(req.ExternalDetails)), &extDetails)
	if err != nil {
		return "", errors.Wrap(err, "Failed to unmarshal ExternalDetails of the request", logan.F{
			"raw_external_details": req.ExternalDetails,
		})
	}

	expectedReference := s.offchainHelper.BuildReference(extDetails.BlockNumber, extDetails.TXHash, offchainAddress, extDetails.OutIndex, 64)
	if expectedReference != string(op.Reference) {
		return fmt.Sprintf("Invalid reference - expected (%s), got (%s).", expectedReference, op.Reference), nil
	}

	if extDetails.Price != amount.One {
		return fmt.Sprintf("Wrong Price - expected (%d), got (%d).", amount.One, extDetails.Price), nil
	}

	lastKnownBlockNumber, err := s.offchainHelper.GetLastKnownBlockNumber()
	if err != nil {
		return "", errors.Wrap(err, "Failed to get last known Block number")
	}
	if lastKnownBlockNumber-extDetails.BlockNumber < s.lastBlocksNotWatch {
		return fmt.Sprintf("Too early to process this Deposit, last existing Offchain Block is (%d), we don't watch last (%d) Blocks",
			lastKnownBlockNumber, s.lastBlocksNotWatch), nil
	}

	block, err := s.offchainHelper.GetBlock(extDetails.BlockNumber)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get Block")
	}

	return s.validateOffchainBlock(*block, extDetails, uint64(req.Amount), offchainAddress), nil
}

func (s *Service) getOffchainAddress(accountAddress, assetName string) (string, error) {
	account, err := s.horizon.Accounts().ByAddress(accountAddress)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get Account by Address")
	}

	for _, extSysAccount := range account.ExternalSystemAccounts {
		if extSysAccount.Type.Name == assetName {
			return extSysAccount.Address, nil
		}
	}

	return "", errNoExtAccount
}

func (s *Service) validateOffchainBlock(block deposit.Block, opExtDetails deposit.ExternalDetails, emissionAmount uint64, offchainAddress string) (checkErr string) {
	for _, tx := range block.TXs {
		if tx.Hash == opExtDetails.TXHash {
			// TX Exists in Offchain
			if len(tx.Outs) <= int(opExtDetails.OutIndex) {
				return fmt.Sprintf("OutIndex is invalid, the Offchain TX has only (%d) Outputs.", len(tx.Outs))
			}
			out := tx.Outs[opExtDetails.OutIndex]

			checkErr = s.validateOffchainOut(out, emissionAmount, offchainAddress)
			if checkErr != "" {
				return checkErr
			}

			return ""
		}
	}

	return fmt.Sprintf("TX with hash (%s) wasn't found in the Block (%d).", opExtDetails.TXHash, opExtDetails.BlockNumber)
}

func (s *Service) validateOffchainOut(out deposit.Out, emissionAmount uint64, offchainAddress string) (checkErr string) {
	if out.Address != offchainAddress {
		return fmt.Sprintf("Invalid Output Address, expected (%s), got (%s).", offchainAddress, out.Address)
	}

	if out.Value < s.offchainHelper.GetMinDepositAmount() {
		return fmt.Sprintf("Output value is less than MinDepositAmount (%d) - using offchain precision.", s.offchainHelper.GetMinDepositAmount())
	}

	valueWithoutFee := out.Value - s.offchainHelper.GetFixedDepositFee()
	systemValue := s.offchainHelper.ConvertToSystem(valueWithoutFee)
	if systemValue != emissionAmount {
		return fmt.Sprintf("Invalid Emission amount, expected (%d), got (%d) - using system precision.", emissionAmount, systemValue)
	}

	return ""
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
