package kyc

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// ParsingData describes the structure of KYC blob retrieved form Horizon.
type parsingData struct {
	FirstName  string      `json:"first_name"`
	LastName   string      `json:"last_name"`
	Address    Address     `json:"address"`
	ETHAddress string      `json:"eth_address"`
	Documents  DocumentsV1 `json:"documents"`

	Version string `json:"version"`
	V2      Data   `json:"v2"`
}

type DocumentsV1 struct {
	KYCIdDocument     string `json:"kyc_id_document"`
	KYCProofOfAddress string `json:"kyc_poa"`
}

func ParseKYCData(data string) (*Data, error) {
	var parsingKYCData parsingData
	err := json.Unmarshal([]byte(data), &parsingKYCData)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal data bytes into Data structure")
	}

	switch parsingKYCData.Version {
	case "v2":
		return &parsingKYCData.V2, nil
	default:
		// v1
		return &Data{
			FirstName:  parsingKYCData.FirstName,
			LastName:   parsingKYCData.LastName,
			Address:    parsingKYCData.Address,
			Documents: Documents{
				IDDocument: IDDocument{
					FaceDocID: parsingKYCData.Documents.KYCIdDocument,
					BackDocID: "",
					Type:      PassportDocType,
				},
				ProofOfAddr: ProofOfAddrDoc{
					FaceFileID: parsingKYCData.Documents.KYCProofOfAddress,
				},
			},
		}, nil
	}
}
