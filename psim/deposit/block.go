package deposit

import "time"

type Block struct {
	Hash      string
	Timestamp time.Time
	TXs       []Tx
}

func (b Block) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"hash":         b.Hash,
		"transactions": len(b.TXs),
	}
}

type Tx struct {
	Hash string
	Outs []Out
}

func (tx Tx) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"hash": tx.Hash,
	}
}

type Out struct {
	Address string
	Value   uint64
}

func (out Out) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"addr":  out.Address,
		"value": out.Value,
	}
}
