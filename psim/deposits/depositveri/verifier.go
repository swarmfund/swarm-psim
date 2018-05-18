package depositveri

import (
	"encoding/json"

	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/deposits/deposit"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
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

// This method is to implement Verifier interface from package verifier.
func (v *Verifier) Run(ctx context.Context) {
	<-ctx.Done()
}

func (v *Verifier) GetOperationType() xdr.OperationType {
	return xdr.OperationTypeCreateIssuanceRequest
}

func (v *Verifier) VerifyOperation(envelope xdr.TransactionEnvelope) (verifyErr, err error) {
	op := envelope.Tx.Operations[0].Body.CreateIssuanceRequestOp

	if op == nil {
		return errors.Errorf("CreateIssuanceRequestOp is nil."), nil
	}

	checkErr, err := v.validateIssuanceOp(*op)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to validate Issuance Op")
	}

	return checkErr, nil
}

// TODO Try make me smaller
func (v *Verifier) validateIssuanceOp(op xdr.CreateIssuanceRequestOp) (verifyErr, err error) {
	req := op.Request

	if string(req.Asset) != v.offchainHelper.GetAsset() {
		return errors.Errorf("Invalid asset - expected (%s), got (%s).", v.offchainHelper.GetAsset(), req.Asset), nil
	}

	accountID, err := v.horizon.Balances().AccountID(req.Receiver.AsString())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get AccountID by Balance from Horizon")
	}
	if accountID == nil {
		return errors.Errorf("No Account was found by provided BalanceID"), nil
	}

	offchainAddress, err := v.getOffchainAddress(*accountID, v.externalSystem)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Offchain Address of the Account")
	}

	extDetails := deposit.ExternalDetails{}
	err = json.Unmarshal([]byte(string(req.ExternalDetails)), &extDetails)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal ExternalDetails of the request", logan.F{
			"raw_external_details": req.ExternalDetails,
		})
	}

	expectedReference := v.offchainHelper.BuildReference(extDetails.BlockNumber, extDetails.TXHash, offchainAddress, extDetails.OutIndex, 64)
	if expectedReference != string(op.Reference) {
		return errors.Errorf("Invalid reference - expected (%s), got (%s).", expectedReference, op.Reference), nil
	}

	if extDetails.Price != amount.One {
		return errors.Errorf("Wrong Price - expected (%d), got (%d).", amount.One, extDetails.Price), nil
	}

	lastKnownBlockNumber, err := v.offchainHelper.GetLastKnownBlockNumber()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get last known Block number")
	}
	if lastKnownBlockNumber-extDetails.BlockNumber < v.lastBlocksNotWatch {
		return errors.Errorf("Too early to process this Deposit, last existing Offchain Block is (%d), we don't watch last (%d) Blocks",
			lastKnownBlockNumber, v.lastBlocksNotWatch), nil
	}

	block, err := v.offchainHelper.GetBlock(extDetails.BlockNumber)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Block")
	}

	return v.verifyOffchainBlock(*block, extDetails, uint64(req.Amount), offchainAddress), nil
}

func (v *Verifier) getOffchainAddress(accountAddress, assetName string) (string, error) {
	account, err := v.horizon.Accounts().ByAddress(accountAddress)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get Account by Address")
	}

	for _, extSysAccount := range account.ExternalSystemAccounts {
		if extSysAccount.AssetCode == assetName {
			return extSysAccount.Address, nil
		}
	}

	return "", errNoExtAccount
}

func (v *Verifier) verifyOffchainBlock(block deposit.Block, opExtDetails deposit.ExternalDetails, emissionAmount uint64, offchainAddress string) (checkErr error) {
	for _, tx := range block.TXs {
		if tx.Hash == opExtDetails.TXHash {
			// TX Exists in Offchain
			if len(tx.Outs) <= int(opExtDetails.OutIndex) {
				return errors.Errorf("OutIndex is invalid, the Offchain TX has only (%d) Outputs.", len(tx.Outs))
			}
			out := tx.Outs[opExtDetails.OutIndex]

			checkErr = v.validateOffchainOut(out, emissionAmount, offchainAddress)
			if checkErr != nil {
				return checkErr
			}

			return nil
		}
	}

	return errors.Errorf("TX with hash (%s) wasn't found in the Block (%d).", opExtDetails.TXHash, opExtDetails.BlockNumber)
}

func (v *Verifier) validateOffchainOut(out deposit.Out, emissionAmount uint64, offchainAddress string) (checkErr error) {
	if out.Address != offchainAddress {
		return errors.Errorf("Invalid Output Address, expected (%s), got (%s).", offchainAddress, out.Address)
	}

	if out.Value < v.offchainHelper.GetMinDepositAmount() {
		return errors.Errorf("Output value is less than MinDepositAmount (%d) - using offchain precision.", v.offchainHelper.GetMinDepositAmount())
	}

	valueWithoutFee := out.Value - v.offchainHelper.GetFixedDepositFee()
	systemValue := v.offchainHelper.ConvertToSystem(valueWithoutFee)
	if systemValue != emissionAmount {
		return errors.Errorf("Invalid Emission amount, expected (%d), got (%d) - using system precision.", emissionAmount, systemValue)
	}

	return nil
}
