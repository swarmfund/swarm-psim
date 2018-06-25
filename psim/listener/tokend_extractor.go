package listener

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/listener/internal"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
)

// TokendExtractor is responsible for taking Txs and extract data from them to be used by processors
type TokendExtractor struct {
	logger *logan.Entry
	Source <-chan horizon.TXPacket
}

// NewTokendExtractor constructs a TokendExtractor using provided source
func NewTokendExtractor(logger *logan.Entry, source <-chan horizon.TXPacket) *TokendExtractor {
	return &TokendExtractor{
		logger: logger,
		Source: source,
	}
}

// TxData holds all data to be processed for broadcasting events
type TxData struct {
	Operations    []xdr.Operation
	OpsResults    []xdr.OperationResult
	LedgerChanges [][]xdr.LedgerEntryChange
	Time          *time.Time
	SourceAccount xdr.AccountId
	PagingToken   string
}

func validateTx(extractedTx horizon.TXPacket) (*TxData, error) {
	extractedTxBody, err := extractedTx.Unwrap()
	if err != nil {
		return nil, errors.Wrap(err, "failed to unwrap tx")
	}

	tx := extractedTxBody.Transaction

	if tx == nil {
		return nil, nil
	}

	txEnvelope, err := tx.SafeEnvelope()

	if err != nil {
		return nil, errors.Wrap(err, "failed to get tx envelope")
	}

	txEnvelopeBody := txEnvelope.Tx

	txLedgerChanges, err := tx.GroupedLedgerChanges()

	if err != nil {
		return nil, errors.Wrap(err, "failed to get grouped ledger changes", logan.F{
			"tx_source_account": txEnvelopeBody.SourceAccount,
		})
	}

	operations := txEnvelopeBody.Operations

	txResult, err := tx.Result()

	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal results", logan.F{
			"tx_source_account": txEnvelopeBody.SourceAccount,
		})
	}

	opsResults, ok := txResult.Result.GetResults()

	if !ok {
		return nil, errors.Wrap(err, "failed to get results", logan.F{
			"tx_source_account": txEnvelopeBody.SourceAccount,
		})
	}

	txTime := &tx.CreatedAt

	pagingToken := tx.PagingToken

	return &TxData{operations, opsResults, txLedgerChanges, txTime, txEnvelopeBody.SourceAccount, pagingToken}, nil
}

// Extract safely gathers all the stuff from each tx and puts it as ExtractedItem to a channel
func (extractor TokendExtractor) Extract(ctx context.Context) <-chan ExtractedItem {
	out := make(chan ExtractedItem)

	go func(chan ExtractedItem) {
		defer func() {
			if r := recover(); r != nil {
				extractor.logger.WithRecover(r).Warn("panic while extracting txdata	")
			}
			close(out)
		}()

		for extractedTx := range extractor.Source {
			select {
			case <-ctx.Done():
				return
			default:
			}

			txData, err := validateTx(extractedTx)
			if err != nil {
				extractor.logger.WithError(err).Warn("got invalid tx")
			}
			if txData == nil {
				continue
			}

			for constructedData := range constructOpData(txData) {
				out <- constructedData
			}
		}
	}(out)

	return out
}

func constructOpData(txData *TxData) <-chan ExtractedItem {
	operations := txData.Operations
	txSourceAccount := txData.SourceAccount
	opsResults := txData.OpsResults
	txLedgerChanges := txData.LedgerChanges
	txTime := txData.Time
	pagingToken := txData.PagingToken

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
				out <- internal.InvalidExtractedItem(errors.New("failed to get tx result tr", logan.F{
					"source_account": txSourceAccount,
				}))
				continue
			}

			opLedgerChanges := txLedgerChanges[currentOpIndex]

			out <- internal.ValidExtractedItem(currentOp, sourceAccount, opLedgerChanges, *opResultTr, txTime, pagingToken)
		}
	}(out)

	return out
}
