package txhandler

import (
	"time"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type txHandler struct {
	statable Statable

	log *logan.Entry
}

func newTxHandler(statable Statable, log *logan.Entry) *txHandler {
	return &txHandler{
		statable: statable,
		log:      log,
	}
}

func (h *txHandler) Handle(tx horizon.Transaction) error {
	envelope := xdr.TransactionEnvelope{}
	err := xdr.SafeUnmarshalBase64(tx.EnvelopeXDR, &envelope)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal tx envelope")
	}

	for _, op := range envelope.Tx.Operations {
		opBody := op.Body
		if opBody.Type != xdr.OperationTypeSetFees {
			continue
		}

		setFeesOp := opBody.MustSetFeesOp()
		if setFeesOp.PayoutsPeriod == nil {
			continue
		}

		payoutPeriod := time.Duration(*setFeesOp.PayoutsPeriod) * time.Second
		if payoutPeriod < 0 {
			// in case of overflow. We should set payout period back to undefined
			h.statable.SetPayoutPeriod(nil)
		} else {
			h.statable.SetPayoutPeriod(&payoutPeriod)
		}

		h.log.WithField("period", payoutPeriod.String()).Info("payout period updated")
	}

	return nil
}
