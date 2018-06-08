package listener

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/listener/internal"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
)

var (
	ErrNilTx = errors.New("empty tx received")
)

// TokendExtractor is responsible for taking Txs and extract data from them to be used by processors
type TokendExtractor <-chan horizon.TXPacket

// TxData holds all data to be processed for broadcasting events
type TxData struct {
	Operations    []xdr.Operation
	OpsResults    []xdr.OperationResult
	LedgerChanges [][]xdr.LedgerEntryChange
	Time          time.Time
	SourceAccount xdr.AccountId
}

func validateTx(extractedTx horizon.TXPacket) (TxData, error) {
	extractedTxBody, err := extractedTx.Unwrap()
	if err != nil {
		return TxData{}, errors.Wrap(err, "failed to unwrap tx")
	}

	tx := extractedTxBody.Transaction

	if tx == nil {
		return TxData{}, ErrNilTx
	}

	txEnvelope, err := tx.SafeEnvelope()

	if err != nil {
		return TxData{}, errors.Wrap(err, "failed to get tx envelope")
	}

	txEnvelopeBody := txEnvelope.Tx

	txLedgerChanges, err := tx.GroupedLedgerChanges()

	if err != nil {
		return TxData{}, errors.Wrap(err, "failed to get grouped ledger changes")
	}

	operations := txEnvelopeBody.Operations

	txResult, err := tx.Result()

	if err != nil {
		return TxData{}, errors.Wrap(err, "failed to unmarshal results")
	}

	opsResults, ok := txResult.Result.GetResults()

	if !ok {
		return TxData{}, errors.Wrap(err, "failed to get results")
	}

	txTime := tx.CreatedAt

	return TxData{operations, opsResults, txLedgerChanges, txTime, txEnvelopeBody.SourceAccount}, nil
}

// Extract safely gathers all the stuff from each tx and puts it as ExtractedItem to a channel
func (extractor TokendExtractor) Extract(ctx context.Context) <-chan ExtractedItem {
	out := make(chan ExtractedItem)

	go func(chan ExtractedItem) {
		defer func() {
			// TODO recover
			close(out)
		}()

		for extractedTx := range extractor {
			select {
			case <-ctx.Done():
				return
			default:
			}

			txData, err := validateTx(extractedTx)
			if err == errors.Cause(ErrNilTx) {
				continue
			}

			if err != nil {
				out <- internal.InvalidExtractedItem(errors.Wrap(err, "invalid tx"))
			}

			for constructedData := range constructOpData(txData) {
				out <- constructedData
			}
		}
	}(out)

	return out
}

func constructOpData(txData TxData) <-chan ExtractedItem {

	operations := txData.Operations
	txSourceAccount := txData.SourceAccount
	opsResults := txData.OpsResults
	txLedgerChanges := txData.LedgerChanges
	txTime := txData.Time

	out := make(chan ExtractedItem)

	go func(chan ExtractedItem) {
		defer func() {
			close(out)
		}()

		sourceAccount := txSourceAccount

		for currentOpIndex, currentOp := range operations {

			if currentOp.SourceAccount != nil {
				sourceAccount = *currentOp.SourceAccount
			}

			opResultTr := opsResults[currentOpIndex].Tr
			if opResultTr == nil {
				out <- internal.InvalidExtractedItem(errors.New("failed to get tx envelope"))
				continue
			}

			opLedgerChanges := txLedgerChanges[currentOpIndex]

			out <- internal.ValidExtractedItem(currentOp, sourceAccount, opLedgerChanges, *opResultTr, txTime)
		}
	}(out)

	return out
}
