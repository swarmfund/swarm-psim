package kyc

// Address is only a nested structure in Data structure.
type Address struct {
	Line1      string `json:"line_1"`
	Line2      string `json:"line_2"`
	City       string `json:"city"` // Use Detroit on Sandbox to receive failed result
	Country    string `json:"country"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
}

func (a Address) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"line_1":      a.Line1,
		"line_2":      a.Line2,
		"city":        a.City,
		"country":     a.Country,
		"state":       a.State,
		"postal_code": a.PostalCode,
	}
}
