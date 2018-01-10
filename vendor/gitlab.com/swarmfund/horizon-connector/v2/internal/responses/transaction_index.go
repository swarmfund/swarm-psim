package responses

import "gitlab.com/swarmfund/horizon-connector/v2/internal/resources"

type TransactionIndex struct {
	Embedded struct {
		Meta    resources.PageMeta      `json:"meta"`
		Records []resources.Transaction `json:"records"`
	} `json:"_embedded"`
}
