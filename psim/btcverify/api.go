package btcverify

import (
	"gitlab.com/tokend/psim/ape"
	"encoding/json"
	"gitlab.com/tokend/psim/ape/problems"
	"fmt"
	"gitlab.com/tokend/psim/psim/bitcoin"
	"github.com/piotrnar/gocoin/lib/btc"
	"net/http"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/go/xdr"
)

func (s *Service) serveAPI() {
	r := ape.DefaultRouter()

	r.Post("/", s.verifyHandler)
	if s.config.Pprof {
		s.log.Info("enabling debugging endpoints")
		ape.InjectPprof(r)
	}

	s.log.WithField("address", s.listener.Addr().String()).Info("listening")

	err := ape.ListenAndServe(s.ctx, s.listener, r)
	if err != nil {
		s.errors <- err
		return
	}
	return
}

// TODO split me to several methods
func (s *Service) verifyHandler(w http.ResponseWriter, r *http.Request) {
	payload := VerifyRequest{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		ape.RenderErr(w, r, problems.BadRequest("Cannot parse JSON request."))
		return
	}

	horizonTX := s.horizon.Transaction(&horizon.TransactionBuilder{
		Envelope: payload.Envelope,
	})

	if len(horizonTX.Operations) != 1 {
		ape.RenderErr(w, r, problems.BadRequest("Provided Transaction envelope contains more than one operation."))
		return
	}

	op := horizonTX.Operations[0]
	if op.Body.Type != xdr.OperationTypeManageCoinsEmissionRequest {
		ape.RenderErr(w, r, problems.BadRequest(fmt.Sprintf(
			"Expected Operation of type ManageCoinEmissionRequest(%d), but got (%d).",
			xdr.OperationTypeManageCoinsEmissionRequest, op.Body.Type)))
		return
	}

	opBody := horizonTX.Operations[0].Body.ManageCoinsEmissionRequestOp

	if opBody.Asset != bitcoin.Asset {
		ape.RenderErr(w, r, problems.BadRequest(fmt.Sprintf("Expected asset to be '%s', but got '%s'.",
			bitcoin.Asset, opBody.Asset)))
		return
	}

	if opBody.Action != xdr.ManageCoinsEmissionRequestActionManageCoinsEmissionRequestCreate {
		ape.RenderErr(w, r, problems.BadRequest("Expected Action to be CER create."))
		return
	}

	reference := bitcoin.BuildCoinEmissionRequestReference(payload.TXHash, payload.OutIndex)
	if string(opBody.Reference) != reference {
		ape.RenderErr(w, r, problems.Conflict(fmt.Sprintf("Expected reference to be '%s', but got '%s'.",
			reference, string(opBody.Reference))))
		return
	}

	block, err := s.btcClient.GetBlockByHash(payload.BlockHash)
	if err != nil {
		s.log.WithField("block_hash", payload.BlockHash).WithError(err).Error("Failed to get Block by hash")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}

	for _, tx := range block.Txs {
		if tx.Hash.String() == payload.TXHash {
			for i, _ := range tx.TxOut {
				if i == payload.OutIndex {
					s.verifyTxOUT(w, r, tx.TxOut[i], opBody, horizonTX)
					return
				}
			}

			p := problems.NotFound("Such TX output was not found in the transaction.")
			(*p.Meta)["tx_hash"] = payload.TXHash
			ape.RenderErr(w, r, p)
			return
		}
	}

	p := problems.NotFound("Such TX was not found in the Block.")
	(*p.Meta)["block_hash"] = payload.BlockHash
	ape.RenderErr(w, r, p)
	return

}

func (s *Service) verifyTxOUT(w http.ResponseWriter, r *http.Request, txOut *btc.TxOut,
	opBody *xdr.ManageCoinsEmissionRequestOp, horizonTX *horizon.TransactionBuilder) {

	// TODO Make sure amount is valid (maybe need to * or / by 10^N)
	if int64(opBody.Amount) != int64(txOut.Value) {
		ape.RenderErr(w, r, problems.Conflict(fmt.Sprintf("Got %d satoshi in the TXOut, but CoinEmissionRequest Amount is %d",
			int64(txOut.Value), int64(opBody.Amount))))
		return
	}

	receiver := opBody.Receiver.AsString()
	account, err := s.horizon.AccountSigned(s.config.Signer, receiver)
	if err != nil {
		s.log.WithField("account_id", receiver).WithError(err).Error("Failed to get Account from Horizon")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}
	if account == nil {
		s.log.WithField("account_id", receiver).Warn("Account from CoinEmissionRequest was not found in Horizon.")
		p := problems.NotFound("Account was not found in Horizon.")
		(*p.Meta)["account_id"] = receiver
		ape.RenderErr(w, r, p)
		return
	}

	// TODO after changing of xdr
	//outAddr := btc.NewAddrFromPkScript(txOut.Pk_script, s.btcClient.IsTestnet())
	//if outAddr != account.BTCAddress {
	//	ape.RenderErr(w, r, problems.Conflict(fmt.Sprintf(
	//		"Was requested CoinEmissionRequest to the '%s' Account, " +
	//			"but its BTC Address does not match the Address from BTC TxOut (%s)", receiver, outAddr)))
	//	return
	//}

	err = horizonTX.Sign(s.config.Signer).Submit()
	if err != nil {
		entry := s.log.WithError(err)

		if serr, ok := err.(horizon.SubmitError); ok {
			opCodes := serr.OperationCodes()
			entry = entry.
				WithField("tx_code", serr.TransactionCode()).
				WithField("op_codes", serr.OperationCodes())
			if len(opCodes) == 1 {
				switch opCodes[0] {
				// safe to move on errors
				case "op_balance_not_found", "reference_duplication":
					entry.Info("tx failed")
					return
				case "op_counterparty_wrong_type":
					entry.Error("tx failed")
					return
				}
			}
		}

		entry.Error("tx failed")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}
}
