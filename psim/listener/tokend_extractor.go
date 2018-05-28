package listener

import (
	"context"

	"gitlab.com/tokend/horizon-connector"
)

type TokendExtractor <-chan horizon.TXPacket

func (extractor TokendExtractor) Extract(ctx context.Context) (<-chan TxData, error) {
	out := make(chan TxData)
	go func() {
		defer func() {
			close(out)
		}()
		for extractedTx := range extractor {
			select {
			case <-ctx.Done():
				return
			default:
				break
			}

			// TODO panic handling
			// TODO reduce cyclomatic complexity
			extractedTxBody, err := extractedTx.Unwrap()
			if err != nil {
				// TODO report error
				continue
			}

			tx := extractedTxBody.Transaction

			if extractedTxBody.Transaction == nil {
				// TODO report error
				continue
			}

			txEnv := tx.Envelope().Tx
			txSourceAccount := txEnv.SourceAccount
			txLedgerChanges := tx.GroupedLedgerChanges()
			ops := txEnv.Operations
			opsResults := tx.Result().Result.MustResults()
			txTime := tx.CreatedAt

			for currentOpIndex, currentOp := range ops {
				opLedgerChanges := txLedgerChanges[currentOpIndex]
				sourceAccount := txSourceAccount

				opResultTr := opsResults[currentOpIndex].Tr
				if opResultTr == nil {
					// TODO report error
					continue
				}
				if currentOp.SourceAccount != nil {
					sourceAccount = *currentOp.SourceAccount
				}

				// TODO nil checking
				out <- TxData{Op: currentOp, SourceAccount: sourceAccount, OpLedgerChanges: opLedgerChanges, OpResult: *opResultTr, CreatedAt: &txTime}
			}
		}
	}()
	return out, nil
}
