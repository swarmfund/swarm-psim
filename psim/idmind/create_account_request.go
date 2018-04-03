package idmind

import (
	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/psim/psim/kyc"
)

// CreateAccountRequest describes the structure of CreateAccount request to IdentityMind.
type CreateAccountRequest struct {
	AccountName string `json:"man"`           // The only required field in the request // 60 chars max
	TxID        string `json:"tid,omitempty"` // 32 chars max

	Email         string `json:"tea"`           // 60 chars max
	FirstName     string `json:"bfn"`           // Max 30 chars
	LastName      string `json:"bln"`           // Max 50 chars
	StreetAddress string `json:"bsn"`           // Max 100 chars
	Country       string `json:"bco,omitempty"` // Max 2 chars ISO 3166-1 alpha-2
	PostalCode    string `json:"bz"`            // Max 20 chars
	City          string `json:"bc"`            // Max 30 chars
	State         string `json:"bs"`            // Max 30 chars Use official postal state/region abbreviations whenever possible (e.g. CA for California)
	//PhoneNumber   string `json:"bc,omitempty"`  // Max 60 chars
	//DateOfBirth string `json:"dob,omitempty"`
}

func (r CreateAccountRequest) validate() error {
	if len(r.AccountName) > 60 {
		return errors.Errorf("AccountName cannot be larger than 60 letters (%s).", r.AccountName)
	}
	if len(r.Email) > 60 {
		return errors.Errorf("Email cannot be larger than 60 letters (%s).", r.Email)
	}
	if len(r.FirstName) > 30 {
		return errors.Errorf("FirstName cannot be larger than 30 letters (%s).", r.FirstName)
	}
	if len(r.LastName) > 50 {
		return errors.Errorf("LastName cannot be larger than 50 letters (%s).", r.LastName)
	}
	if len(r.StreetAddress) > 100 {
		return errors.Errorf("StreetAddress cannot be larger than 100 letters (%s).", r.StreetAddress)
	}
	if len(r.Country) > 2 {
		return errors.Errorf("Country cannot be larger than 2 letters (%s).", r.Country)
	}
	if len(r.PostalCode) > 20 {
		return errors.Errorf("PostalCode cannot be larger than 20 letters (%s).", r.PostalCode)
	}
	if len(r.City) > 30 {
		return errors.Errorf("City cannot be larger than 30 letters (%s).", r.City)
	}
	if len(r.State) > 30 {
		return errors.Errorf("State cannot be larger than 30 letters (%s).", r.State)
	}

	return nil
}

func buildCreateAccountRequest(data kyc.Data, email string) (*CreateAccountRequest, error) {
	countryCode, err := convertToISO(data.Address.Country)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to convert Country '%s' to ISO", data.Address.Country))
	}

	r := CreateAccountRequest{
		// Not sure email is a good data to put here
		AccountName: email,

		Email:         email,
		FirstName:     data.FirstName,
		LastName:      data.LastName,
		StreetAddress: fmt.Sprintf("%s %s", data.Address.Line1, data.Address.Line2),
		Country:       countryCode,
		PostalCode:    data.Address.PostalCode,
		City:          data.Address.City,
		State:         data.Address.State,
	}

	validateErr := r.validate()
	if validateErr != nil {
		return nil, validateErr
	}

	return &r, nil
}
