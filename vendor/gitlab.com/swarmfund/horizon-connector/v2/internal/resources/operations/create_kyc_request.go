package operations

type CreateKYCRequest struct {
	RequestID          uint64                 `json:"request_id"`
	AccountToUpdateKYC string                 `json:"account_to_update_kyc"`
	KYCData            map[string]interface{} `json:"kyc_data"`
	PT                 string                 `json:"paging_token"`
	TransactionID      string                 `json:"transaction_id"`
}
