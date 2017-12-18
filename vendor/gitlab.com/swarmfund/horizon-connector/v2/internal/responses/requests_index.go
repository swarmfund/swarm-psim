package responses

import "gitlab.com/swarmfund/horizon-connector/v2/internal/resources"

type RequestsIndex struct {
	Embedded struct {
		Records []resources.Request `json:"records"`
	} `json:"_embedded"`
}
