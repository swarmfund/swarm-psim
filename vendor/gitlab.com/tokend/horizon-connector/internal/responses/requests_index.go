package responses

import "gitlab.com/tokend/horizon-connector/internal/resources"

type RequestsIndex struct {
	Embedded struct {
		Records []resources.Request `json:"records"`
	} `json:"_embedded"`
}
