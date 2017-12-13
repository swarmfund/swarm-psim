package create_account_streamer

type OperationMap map[string]interface{}

type OperationsResponse struct {
	Embedded struct {
		Records []map[string]interface{}
	} `json:"_embedded"`

	Links struct {
		Next struct {
			HREF string `json:"href"`
		} `json:"next"`
	} `json:"_links"`
}
