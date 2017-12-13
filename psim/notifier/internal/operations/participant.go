package operations

type Participant struct {
	AccountID string      `json:"account_id,omitempty"`
	BalanceID string      `json:"balance_id,omitempty"`
	Email     string      `json:"email,omitempty"`
	Effects   BaseEffects `json:"effects,omitempty"`
}

type ApiParticipant struct {
	AccountID string `json:"account_id,omitempty"`
	BalanceID string `json:"balance_id,omitempty"`
	Email     string `json:"email,omitempty"`
}

func (p *Participant) fromApiParticipant(ap *ApiParticipant) {
	p.Email = ap.Email
}
