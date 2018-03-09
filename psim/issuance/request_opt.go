package issuance

type RequestOpt struct {
	Reference string
	Receiver  string
	Asset     string
	Amount    uint64
	Details   string
}

func (i RequestOpt) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"reference": i.Reference,
		"receiver":  i.Receiver,
		"asset":     i.Asset,
		"amount":    i.Amount,
		"details":   i.Details,
	}
}
