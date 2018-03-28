package idmind

import "fmt"

// CreateAccountRequest describes the structure of CreateAccount request to IdentityMind.
type CreateAccountRequest struct {
	AccountName string `json:"man"`           // The only required field in the request
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

// TODO
func (r CreateAccountRequest) validate() error {
	// TODO Check restrictions on fields length
	return nil
}

func buildCreateAccountRequest(data KYCData, email string) (*CreateAccountRequest, error) {
	r := CreateAccountRequest{
		// FIXME
		AccountName: "hardcoded account name",

		Email:         email,
		FirstName:     data.FirstName,
		LastName:      data.LastName,
		StreetAddress: fmt.Sprintf("%s %s", data.Address.Line1, data.Address.Line2),
		Country:       convertToISO(data.Address.Country),
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
