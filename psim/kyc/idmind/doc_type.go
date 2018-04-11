package idmind

type DocType string

const (
	PassportDocType        DocType = "PP"
	DrivingLicenseDocType  DocType = "DL"
	IdentityCardDocType    DocType = "IC"
	ResidencePermitDocType DocType = "RP"
)

var validDocTypes = map[DocType]struct{}{
	PassportDocType:        {},
	DrivingLicenseDocType:  {},
	IdentityCardDocType:    {},
	ResidencePermitDocType: {},
}
