package ratesync

import (
	"gitlab.com/swarmfund/horizon-connector"
)

type SyncResults struct {
	slice [][]horizon.SetRateOp
	//mux   sync.Mutex
}

func NewSyncResults() *SyncResults {
	return &SyncResults{
		slice: make([][]horizon.SetRateOp, 10),
	}
}

func (r *SyncResults) Set(index int64, ops []horizon.SetRateOp) {
	r.slice[index%int64(len(r.slice))] = ops
}

func (r *SyncResults) Get(index int64) []horizon.SetRateOp {
	return r.slice[index%int64(len(r.slice))]
}

//func (r *SyncResults) Append(v SyncResult) {
//	r.mux.Lock()
//	defer r.mux.Unlock()
//	if len(r.slice) == 10 {
//		r.slice = r.slice[1:]
//	}
//	r.slice = append(r.slice, v)
//}
//
//func (r *SyncResults) Get(sync int64) (SyncResult, bool) {
//	r.mux.Lock()
//	defer r.mux.Unlock()
//	for _, r := range r.slice {
//		if r.Sync == sync {
//			return r, true
//		}
//	}
//	return SyncResult{}, false
//}
//
//func (r *SyncResults) Shift() (v SyncResult, ok bool) {
//	r.mux.Lock()
//	defer r.mux.Unlock()
//	if len(r.slice) == 0 {
//		return
//	}
//	v, r.slice, ok = r.slice[0], r.slice[1:], true
//	return
//}
