package kyc

import "time"

var usaTerritories = []string{
	"United States of America", "US",
	// Just in case
	"USA",
	"United States",

	// https://en.wikipedia.org/wiki/Unincorporated_territories_of_the_United_States
	// https://en.wikipedia.org/wiki/Territories_of_the_United_States#Inhabited_territories_2

	"United States Minor Outlying Islands", "UM",
	"Guam", "GU",
	"Northern Mariana Islands", "MP",
	"Virgin Islands, U.S.", "VI",
	"American Samoa", "AS",
	"Puerto Rico", "PR",
}

type Data struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Address     Address   `json:"address"`
	Documents   Documents `json:"documents"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

func (d Data) IsUSA() bool {
	for _, terr := range usaTerritories {
		if d.Address.Country == terr {
			return true
		}
	}

	return false
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
