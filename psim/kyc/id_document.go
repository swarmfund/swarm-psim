package kyc

type IDDocument struct {
	FaceFile DocFile  `json:"front"`
	BackFile *DocFile `json:"back"`
	Type     DocType  `json:"type"`
}

func (d IDDocument) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"face": d.FaceFile.ID,
		"back": d.BackFile.ID,
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
