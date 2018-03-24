package mrefairdrop

type bonusParams struct {
	IsVerified bool
	BalanceID  string
	Balance    uint64
	Referrals  map[string]struct{}
}

func newBonusParams() bonusParams {
	return bonusParams{
		Referrals: make(map[string]struct{}),
	}
}

func (p *bonusParams) addReferral(accID string) {
	// It actually OK if the accID is already in Referrals, because Referrals is a set.
	p.Referrals[accID] = struct{}{}
}

func (p *bonusParams) deleteReferral(accID string) {
	// It actually OK if the accID does not exist int Referrals.
	delete(p.Referrals, accID)
}
