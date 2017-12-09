package noor

import "gitlab.com/swarmfund/horizon-connector"

type Tick struct {
	ops []horizon.SetRateOp
}

func (t Tick) Ops() []horizon.SetRateOp {
	return t.ops
}

func (t Tick) Clone() Tick {
	var result Tick
	result.ops = make([]horizon.SetRateOp, len(t.ops))
	copy(result.ops, t.ops)
	return result
}
