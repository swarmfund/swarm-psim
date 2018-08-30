package masternode

type candidate struct {
	address      string
	promoteAfter int64
}

type RewardQueue struct {
	i            int
	current      []string
	all          []string
	candidates   []*candidate
	promoteAfter int64
	blacklist    []string
}

func (q *RewardQueue) Add(address string) {
	if q.isKnown(address) || q.isBlacklisted(address) {
		return
	}
	q.candidates = append(q.candidates, &candidate{
		address:      address,
		promoteAfter: q.promoteAfter,
	})
}

func (q *RewardQueue) isBlacklisted(address string) bool {
	for _, v := range q.blacklist {
		if v == address {
			return true
		}
	}
	return false
}

func (q *RewardQueue) isKnown(address string) bool {
	for _, v := range q.all {
		if v == address {
			return true
		}
	}
	for _, v := range q.candidates {
		if v.address == address {
			return true
		}
	}
	return false
}

func (q *RewardQueue) Next() *string {
	// promote candidates
	for _, candidate := range q.candidates {
		if candidate.promoteAfter == 0 {
			q.all = append(q.all, candidate.address)
		}
		candidate.promoteAfter -= 1
	}
	if len(q.all) <= q.i {
		q.i = 0
	}
	if len(q.all) == 0 {
		return nil
	}
	next := q.all[q.i]
	q.i += 1
	return &next
}

func (q *RewardQueue) Remove(address string) {
	for i, v := range q.all {
		if v == address {
			q.all = q.all[:i+copy(q.all[i:], q.all[i+1:])]
		}
	}
	for i, v := range q.candidates {
		if v.address == address {
			q.candidates = q.candidates[:i+copy(q.candidates[i:], q.candidates[i+1:])]
		}
	}
}
