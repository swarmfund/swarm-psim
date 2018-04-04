package kyc

type Data struct {
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Address    Address   `json:"address"`
	ETHAddress string    `json:"eth_address"`
	Documents  Documents `json:"documents"`
	Sequence   string    `json:"sequence"`
}

func (d Data) IsUSA() bool {
	return d.Address.Country == "United States of America" || d.Address.Country == "US" || d.Address.Country == "USA" ||
		d.Address.Country == "United States"
}

func (d Data) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"first_name":  d.FirstName,
		"last_name":   d.LastName,
		"address":     d.Address,
		"eth_address": d.ETHAddress,
		"documents":   d.Documents,
		"sequence":    d.Sequence,
	}
}

type Documents struct {
	IDDocument  IDDocument     `json:"kyc_id_document"`
	ProofOfAddr ProofOfAddrDoc `json:"kyc_poa"`
}

// TODO
//func (d Documents) GetLoganFields() map[string]interface{} {
//	return map[string]interface{}{
//		"kyc_id":               d.KYCIdDocument,
//		"kyc_proof_of_address": d.KYCProofOfAddress,
//	}
//}

type IDDocument struct {
	FaceFileID string  `json:"front"`
	BackFileID string  `json:"back"`
	Type       DocType `json:"type"`
}

type DocType string

const (
	PassportDocType        DocType = "passport"
	DrivingLicenseDocType  DocType = "driving_license"
	IdentityCardDocType    DocType = "identity_card"
	ResidencePermitDocType DocType = "residence_permit"
)

type ProofOfAddrDoc struct {
	Face string `json:"front"`
}
