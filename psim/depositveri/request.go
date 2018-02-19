package depositveri

import "gitlab.com/swarmfund/go/xdr"

type request struct {
	AccountID      string
	Envelope       xdr.TransactionEnvelope
	EnvelopeString string
}

func (r request) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"account_id": r.AccountID,
		"envelope":   r.EnvelopeString,
	}
}

func (r request) GetEnvelope() xdr.TransactionEnvelope {
	return r.Envelope
}
