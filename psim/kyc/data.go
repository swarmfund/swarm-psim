package kyc

import "time"

type Data struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Address     Address   `json:"address"`
	Documents   Documents `json:"documents"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

func (d Data) IsUSA() bool {
	return d.Address.Country == "United States of America" || d.Address.Country == "US" || d.Address.Country == "USA" ||
		d.Address.Country == "United States"
}

func (d Data) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"first_name":    d.FirstName,
		"last_name":     d.LastName,
		"address":       d.Address,
		"documents":     d.Documents,
		"date_of_birth": d.DateOfBirth,
	}
}
