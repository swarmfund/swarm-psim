package bitcoin

// Out describes an Output of a TX.
// Out is used to be passed to BitcoinCore into CreateRawTX request,
// describing the Inputs of a new TX (the previous unspent Outputs).
type Out struct {
	TXHash string `json:"txid"`
	Vout   uint32 `json:"vout"`
}

func (o Out) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"tx_hash": o.TXHash,
		"vout":    o.Vout,
	}
}

// InputUTXO describes the Input of a new TX - a UTXO being spent by this Input.
// InputUTXO is used to be passed to BitcoinCore into SignRawTX request.
type InputUTXO struct {
	Out
	ScriptPubKey string  `json:"scriptPubKey"`
	RedeemScript *string `json:"redeemScript,omitempty"`
}

func (i InputUTXO) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"out":            i.Out,
		"script_pub_key": i.ScriptPubKey,
		"redeem_script":  i.RedeemScript,
	}
}

type WalletUTXO struct {
	InputUTXO

	Address       string  `json:"address"`
	Amount        float64 `json:"amount"`
	Confirmations uint    `json:"confirmations"`
}

func (w WalletUTXO) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"out":            w.Out,
		"script_pub_key": w.ScriptPubKey,
		"redeem_script":  w.RedeemScript,

		"address":       w.Address,
		"amount":        w.Amount,
		"confirmations": w.Confirmations,
	}
}

// UTXO is the structure of UTXO returned by BitcoinCore.
// UTXO doesn't actually contains data about which Output of which TX this UTXO is connected with.
type UTXO struct {
	Confirmations uint         `json:"confirmations"`
	Value         float64      `json:"value"`
	ScriptPubKey  ScriptPubKey `json:"scriptPubKey"`
}

func (u UTXO) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"confirmations":  u.Confirmations,
		"value":          u.Value,
		"script_pub_key": u.ScriptPubKey,
	}
}

// ScriptPubKey is used to store data about ScriptPubKey in the structure of UTXO.
type ScriptPubKey struct {
	Hex       string   `json:"hex"`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses"`
}

func (s ScriptPubKey) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"hex":       s.Hex,
		"type":      s.Type,
		"addresses": s.Addresses,
	}
}

// FundResult is the response from BitcoinCore for the FundRawTX request.
type FundResult struct {
	Hex            string  `json:"hex"`
	ChangePosition int     `json:"changepos"`
	FeePaid        float64 `json:"fee"`
}

func (f FundResult) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"hex":             f.Hex,
		"change_position": f.ChangePosition,
		"fee_paid":        f.FeePaid,
	}
}
