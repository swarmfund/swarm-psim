package kyc

import (
	"encoding/json"

	"github.com/pkg/errors"
	"gitlab.com/tokend/horizon-connector"
)

// ParsingData describes the structure of KYC Blob retrieved form Horizon.
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

func ParseKYCBlob(blob *horizon.Blob) (*Data, error) {

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
			FirstName: parsingKYCData.FirstName,
			LastName:  parsingKYCData.LastName,
			Address:   parsingKYCData.Address,
			Documents: Documents{
				IDDocument: IDDocument{
					FaceFile: DocFile{
						ID: parsingKYCData.Documents.KYCIdDocument,
					},
					BackFile: nil,
					Type:     PassportDocType,
				},
				ProofOfAddr: ProofOfAddrDoc{
					FaceFile: DocFile{
						ID: parsingKYCData.Documents.KYCProofOfAddress,
					},
				},
			},
		}, nil
	}
}

// ParsingFirstNameData describes the structure of shortened(FirstName only) KYC Blob retrieved form Horizon.
type parsingFirstNameData struct {
	FirstName string `json:"first_name"`

	Version string        `json:"version"`
	V2      FirstNameData `json:"v2"`
}

func ParseKYCFirstName(data string) (string, error) {
	var parsingData parsingFirstNameData
	err := json.Unmarshal([]byte(data), &parsingData)
	if err != nil {
		return "", errors.Wrap(err, "Failed to unmarshal data bytes into Data structure")
	}

	switch parsingData.Version {
	case "v2":
		return parsingData.V2.FirstName, nil
	default:
		// v1
		return parsingData.FirstName, nil
	}
}
