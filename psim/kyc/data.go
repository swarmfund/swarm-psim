package kyc

type Data struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Address   Address   `json:"address"`
	Documents Documents `json:"documents"`
}

func (d Data) IsUSA() bool {
	return d.Address.Country == "United States of America" || d.Address.Country == "US" || d.Address.Country == "USA" ||
		d.Address.Country == "United States"
}

func (d Data) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"first_name": d.FirstName,
		"last_name":  d.LastName,
		"address":    d.Address,
		"documents":  d.Documents,
	}
}

type Documents struct {
	IDDocument  IDDocument     `json:"kyc_id_document"`
	ProofOfAddr ProofOfAddrDoc `json:"kyc_poa"`
}

func (d Documents) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"id":               d.IDDocument,
		"proof_of_address": d.ProofOfAddr,
	}
}

type IDDocument struct {
	FaceDocID string  `json:"front"`
	BackDocID string  `json:"back"`
	Type      DocType `json:"type"`
}

func (d IDDocument) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"face": d.FaceDocID,
		"back": d.BackDocID,
		"type": d.Type,
	}
}

type DocType string

const (
	PassportDocType        DocType = "passport"
	DrivingLicenseDocType  DocType = "driving_license"
	IdentityCardDocType    DocType = "identity_card"
	ResidencePermitDocType DocType = "residence_permit"
)

type ProofOfAddrDoc struct {
	FaceFileID string `json:"front"`
}

func (d ProofOfAddrDoc) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"face": d.FaceFileID,
	}
}
