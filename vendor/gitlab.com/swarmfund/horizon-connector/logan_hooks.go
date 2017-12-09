package horizon

import "gitlab.com/distributed_lab/logan"

const (
	opCodesKey = "op-codes"
	txCodeKey  = "tx-code"
)

type TXFailedHook struct{}

func (h *TXFailedHook) Levels() []logan.Level {
	return logan.AllLevels
}

func (h *TXFailedHook) Fire(entry *logan.Entry) error {
	if err, exists := entry.Entry.Data[logan.ErrorKey]; exists {
		if serr, ok := err.(SubmitError); ok {
			entry.Data[opCodesKey] = serr.OperationCodes()
			entry.Data[txCodeKey] = serr.TransactionCode()
		}
	}
	return nil
}
