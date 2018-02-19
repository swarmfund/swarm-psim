package depositveri

import (
	"fmt"

	"encoding/json"

	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/amount"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/ape/problems"
	"gitlab.com/swarmfund/psim/psim/deposit"
	"gitlab.com/swarmfund/psim/psim/verification"
	"gitlab.com/swarmfund/psim/psim/verifier"
)

var errNoExtAccount = errors.New("External system Account was not found.")

type Verifier struct {
	externalSystem     string
	log                *logan.Entry
	lastBlocksNotWatch uint64
	// TODO Interface
	horizon        *horizon.Connector
	offchainHelper deposit.OffchainHelper
}

func newVerifier(
	serviceName string,
	externalSystem string,
	log *logan.Entry,
	lastBlocksNotWatch uint64,
	horizon *horizon.Connector,
	offchainHelper deposit.OffchainHelper) *Verifier {

	return &Verifier{
		externalSystem:     externalSystem,
		log:                log.WithField("service", serviceName),
		lastBlocksNotWatch: lastBlocksNotWatch,
		horizon:            horizon,
		offchainHelper:     offchainHelper,
	}
}

func (v *Verifier) ReadRequest(w http.ResponseWriter, r *http.Request) (verifier.Request, bool) {
	req := deposit.VerifyRequest{}
	if ok := verification.ReadAPIRequest(v.log, w, r, &req); !ok {
		return nil, false
	}

	envelope := xdr.TransactionEnvelope{}
	err := envelope.Scan(req.Envelope)
	if err != nil {
		v.log.WithField("envelope_string", req.Envelope).WithError(err).
			Warn("Failed to Scan TransactionEnvelope from string in request.")
		ape.RenderErr(w, r, problems.BadRequest("Cannot parse Envelope from string."))
		return nil, false
	}

	return request{
		AccountID:      req.AccountID,
		Envelope:       envelope,
		EnvelopeString: req.Envelope,
	}, true
}

func (v *Verifier) VerifyRequest(r verifier.Request) (verifyErr, err error) {
	checkErr, err := v.validateDeposit(r.(request))
	if err != nil {
		return nil, err
	}

	if checkErr != "" {
		return errors.New(checkErr), nil
	}

	return nil, nil
}

func (v *Verifier) validateDeposit(req request) (checkErr string, err error) {
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

	checkErr, err = v.validateIssuanceOp(*op, req.AccountID)
	if err != nil {
		return "", errors.Wrap(err, "Failed to validate Issuance Op")
	}

	return checkErr, nil
}

// TODO Try make me smaller
func (v *Verifier) validateIssuanceOp(op xdr.CreateIssuanceRequestOp, accountAddress string) (checkErr string, err error) {
	req := op.Request

	if string(req.Asset) != v.offchainHelper.GetAsset() {
		return fmt.Sprintf("Invalid asset - expected (%s), got (%s).", v.offchainHelper.GetAsset(), req.Asset), nil
	}

	balanceID, err := deposit.GetBalanceID(v.horizon, accountAddress, v.offchainHelper.GetAsset())
	if err != nil {
		return "", errors.Wrap(err, "Failed to get BalanceID by AccountAddress and Asset from Horizon")
	}
	if req.Receiver.AsString() != balanceID {
		return fmt.Sprintf("Invalid BalanceID - expected (%s), got (%s).", balanceID, req.Receiver.AsString()), nil
	}

	offchainAddress, err := v.getOffchainAddress(accountAddress, v.externalSystem)
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

	expectedReference := v.offchainHelper.BuildReference(extDetails.BlockNumber, extDetails.TXHash, offchainAddress, extDetails.OutIndex, 64)
	if expectedReference != string(op.Reference) {
		return fmt.Sprintf("Invalid reference - expected (%s), got (%s).", expectedReference, op.Reference), nil
	}

	if extDetails.Price != amount.One {
		return fmt.Sprintf("Wrong Price - expected (%d), got (%d).", amount.One, extDetails.Price), nil
	}

	lastKnownBlockNumber, err := v.offchainHelper.GetLastKnownBlockNumber()
	if err != nil {
		return "", errors.Wrap(err, "Failed to get last known Block number")
	}
	if lastKnownBlockNumber-extDetails.BlockNumber < v.lastBlocksNotWatch {
		return fmt.Sprintf("Too early to process this Deposit, last existing Offchain Block is (%d), we don't watch last (%d) Blocks",
			lastKnownBlockNumber, v.lastBlocksNotWatch), nil
	}

	block, err := v.offchainHelper.GetBlock(extDetails.BlockNumber)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get Block")
	}

	return v.validateOffchainBlock(*block, extDetails, uint64(req.Amount), offchainAddress), nil
}

func (v *Verifier) getOffchainAddress(accountAddress, assetName string) (string, error) {
	account, err := v.horizon.Accounts().ByAddress(accountAddress)
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

func (v *Verifier) validateOffchainBlock(block deposit.Block, opExtDetails deposit.ExternalDetails, emissionAmount uint64, offchainAddress string) (checkErr string) {
	for _, tx := range block.TXs {
		if tx.Hash == opExtDetails.TXHash {
			// TX Exists in Offchain
			if len(tx.Outs) <= int(opExtDetails.OutIndex) {
				return fmt.Sprintf("OutIndex is invalid, the Offchain TX has only (%d) Outputs.", len(tx.Outs))
			}
			out := tx.Outs[opExtDetails.OutIndex]

			checkErr = v.validateOffchainOut(out, emissionAmount, offchainAddress)
			if checkErr != "" {
				return checkErr
			}

			return ""
		}
	}

	return fmt.Sprintf("TX with hash (%s) wasn't found in the Block (%d).", opExtDetails.TXHash, opExtDetails.BlockNumber)
}

func (v *Verifier) validateOffchainOut(out deposit.Out, emissionAmount uint64, offchainAddress string) (checkErr string) {
	if out.Address != offchainAddress {
		return fmt.Sprintf("Invalid Output Address, expected (%s), got (%s).", offchainAddress, out.Address)
	}

	if out.Value < v.offchainHelper.GetMinDepositAmount() {
		return fmt.Sprintf("Output value is less than MinDepositAmount (%d) - using offchain precision.", v.offchainHelper.GetMinDepositAmount())
	}

	valueWithoutFee := out.Value - v.offchainHelper.GetFixedDepositFee()
	systemValue := v.offchainHelper.ConvertToSystem(valueWithoutFee)
	if systemValue != emissionAmount {
		return fmt.Sprintf("Invalid Emission amount, expected (%d), got (%d) - using system precision.", emissionAmount, systemValue)
	}

	return ""
}
