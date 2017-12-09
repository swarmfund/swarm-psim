package handlers

import (
	"net/http"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/psim/ape"
	"gitlab.com/tokend/psim/ape/problems"
	"gitlab.com/tokend/psim/psim/taxman/internal/resource"
	"github.com/pkg/errors"
	"gitlab.com/tokend/psim/psim/taxman/internal/snapshoter"
)

func Verify(w http.ResponseWriter, r *http.Request) {
	request := resource.SyncRequest{}
	if err := ape.Bind(r, &request); err != nil {
		ape.RenderErr(w, r, problems.BadRequest("Cannot decode JSON request."))
		return
	}

	snapshot := Snapshots(r).Get(request.Ledger)
	if snapshot == nil {
		ape.RenderErr(w, r, problems.NotFound(""))
		return
	}

	if err := verifySyncRequest(&request, snapshot); err != nil {
		ape.RenderErr(w, r, problems.BadRequest(""))
		return
	}

	for _, txenv := range request.Transactions {
		// TODO decide on multi-sign based on config
		err := Horizon(r).Transaction(&horizon.TransactionBuilder{
			Envelope: txenv,
		}).Sign(Signer(r)).Submit()
		//err := a.service.horizon.SubmitTX(txenv)
		if err != nil {
			if serr, ok := err.(horizon.SubmitError); ok {
				for _, opcode := range serr.OperationCodes() {
					switch opcode {
					case "op_reference_duplication":
						// that's ok, transaction was submitted previously
					default:
						Log(r).
							WithField("tx code", serr.TransactionCode()).
							WithField("op codes", serr.OperationCodes()).
							Warn("failed to submit payout tx")
						//a.service.errors <- errors.Wrap(err, "payout tx failed")
						ape.RenderErr(w, r, problems.ServerError(err))
					}
				}
			} else {
				ape.RenderErr(w, r, problems.ServerError(err))
				return
			}
		}
	}
}

var (
	ErrSyncRequestInvalid = errors.New("sync request is invalid in some way")
)

func verifySyncRequest(request *resource.SyncRequest, snapshot *snapshoter.Snapshot) error {
	// TODO check if op count matches

	for _, txenv := range request.Transactions {
		env := xdr.TransactionEnvelope{}
		err := xdr.SafeUnmarshalBase64(txenv, &env)
		if err != nil {
			return errors.Wrap(err, "tx malformed")
		}
		for _, op := range env.Tx.Operations {
			if op.Body.Type != xdr.OperationTypePayment {
				return ErrSyncRequestInvalid
			}
			body := op.Body.PaymentOp
			snapOp, ok := snapshot.SyncState[string(body.Reference)]
			if !ok {
				return ErrSyncRequestInvalid
			}
			if snapOp.Amount != int64(body.Amount) {
				return ErrSyncRequestInvalid
			}
			if string(snapOp.DestinationBalanceID) != body.DestinationBalanceId.AsString() {
				return ErrSyncRequestInvalid
			}
			if string(snapOp.SourceBalanceID) != body.SourceBalanceId.AsString() {
				return ErrSyncRequestInvalid
			}
		}
	}

	return nil
}
