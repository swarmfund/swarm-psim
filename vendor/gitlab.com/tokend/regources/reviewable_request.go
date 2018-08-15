package regources

import "encoding/json"

// Represents Reviewable request
type ReviewableRequest struct {
	ID           uint64                    `json:"id,string"`
	PT           string                    `json:"paging_token"`
	Requestor    string                    `json:"requestor"`
	Reviewer     string                    `json:"reviewer"`
	Reference    *string                   `json:"reference"`
	RejectReason string                    `json:"reject_reason"`
	Hash         string                    `json:"hash"`
	Details      *ReviewableRequestDetails `json:"details"`
	CreatedAt    Time                      `json:"created_at"`
	UpdatedAt    Time                      `json:"updated_at"`

	// RequestStateI  - integer representation of request state
	State int32 `json:"request_state_i"`
	// RequestState  - string representation of request state
	StateName string `json:"request_state"`
}

// TODO GetLoganFields implementation

func (r *ReviewableRequest) PagingToken() string {
	return r.PT
}

// ReviewableRequestDetails - provides specific for request type details.
// Note: json key of specific request must be equal to xdr.ReviewableRequestType.ShortString result
type ReviewableRequestDetails struct {
	// RequestType  - integer representation of request type
	RequestType int32 `json:"request_type_i"`
	// RequestTypeName  - string representation of request type
	RequestTypeName string `json:"request_type"`

	AssetCreate       *AssetCreationRequest `json:"asset_create,omitempty"`
	AssetUpdate       *AssetUpdateRequest   `json:"asset_update,omitempty"`
	PreIssuanceCreate *PreIssuanceRequest   `json:"pre_issuance_create,omitempty"`
	IssuanceCreate    *IssuanceRequest      `json:"issuance_create,omitempty"`
	Withdraw          *WithdrawalRequest    `json:"withdraw,omitempty"`
	TwoStepWithdraw   *WithdrawalRequest    `json:"two_step_withdrawal"`
	Sale              *SaleCreationRequest  `json:"sale,omitempty"`
	LimitsUpdate           *LimitsUpdateRequest      `json:"limits_update,omitempty"`
	AMLAlert               *AMLAlertRequest          `json:"aml_alert"`
	KYC                    *UpdateKYCRequest         `json:"update_kyc,omitempty"`
	UpdateSaleDetails      *UpdateSaleDetailsRequest `json:"update_sale_details"`
	UpdateSaleEndTime      *UpdateSaleEndTimeRequest `json:"update_sale_end_time"`
	PromotionUpdateRequest *PromotionUpdateRequest   `json:"promotion_update_request"`
	Invoice                *InvoiceRequest           `json:"invoice,omitempty"`
	Contract               *ContractRequest          `json:"contract,omitempty"`
}

type AMLAlertRequest struct {
	BalanceID string `json:"balance_id"`
	Amount    Amount `json:"amount"`
	Reason    string `json:"reason"`
}

type AssetCreationRequest struct {
	Code                   string                 `json:"code"`
	Policies               []Flag                 `json:"policies"`
	PreIssuedAssetSigner   string                 `json:"pre_issued_asset_signer"`
	MaxIssuanceAmount      Amount                 `json:"max_issuance_amount"`
	InitialPreissuedAmount Amount                 `json:"initial_preissued_amount"`
	Details                map[string]interface{} `json:"details"`
}

type AssetUpdateRequest struct {
	Code     string                 `json:"code"`
	Policies []Flag                 `json:"policies"`
	Details  map[string]interface{} `json:"details"`
}

type ContractRequest struct {
	Escrow    string                 `json:"escrow"`
	Details   map[string]interface{} `json:"details"`
	StartTime Time                   `json:"start_time"`
	EndTime   Time                   `json:"end_time"`
}

type InvoiceRequest struct {
	Amount     Amount                 `json:"amount"`
	Asset      string                 `json:"asset"`
	ContractID string                 `json:"contract_id,omitempty"`
	Details    map[string]interface{} `json:"details"`
}

type IssuanceRequest struct {
	Asset           string                 `json:"asset"`
	Amount          Amount                 `json:"amount"`
	Receiver        string                 `json:"receiver"`
	ExternalDetails map[string]interface{} `json:"external_details"`
}

type LimitsUpdateRequest struct {
	DocumentHash string                 `json:"document_hash"`
	Details      map[string]interface{} `json:"details"`
}

type PreIssuanceRequest struct {
	Asset     string `json:"asset"`
	Amount    Amount `json:"amount"`
	Signature string `json:"signature"`
	Reference string `json:"reference"`
}

type PromotionUpdateRequest struct {
	SaleID           uint64              `json:"sale_id"`
	NewPromotionData SaleCreationRequest `json:"new_promotion_data"`
}

type SaleCreationRequest struct {
	BaseAsset           string                 `json:"base_asset"`
	DefaultQuoteAsset   string                 `json:"default_quote_asset"`
	StartTime           Time                   `json:"start_time"`
	EndTime             Time                   `json:"end_time"`
	SoftCap             string                 `json:"soft_cap"`
	HardCap             string                 `json:"hard_cap"`
	SaleType            Flag                   `json:"sale_type"`
	BaseAssetForHardCap string                 `json:"base_asset_for_hard_cap"`
	Details             map[string]interface{} `json:"details"`
	QuoteAssets         []SaleQuoteAsset       `json:"quote_assets"`
	State               Flag                   `json:"state"`
}

type SaleQuoteAsset struct {
	QuoteAsset string `json:"quote_asset"`
	Price      Amount `json:"price"`
}

type UpdateKYCRequest struct {
	AccountToUpdateKYC string                 `json:"account_to_update_kyc"`
	AccountTypeToSet   AccountTypeToSet       `json:"account_type_to_set"`
	KYCLevel           uint32                 `json:"kyc_level"`
	KYCData            map[string]interface{} `json:"kyc_data"`
	// KYCDataStruct is the data from raw map of KYCData, unmarshalled into typed struct in custom Unmarshal below
	KYCDataStruct   KYCData                  `json:"-"`
	AllTasks        uint32                   `json:"all_tasks"`
	PendingTasks    uint32                   `json:"pending_tasks"`
	SequenceNumber  uint32                   `json:"sequence_number"`
	ExternalDetails []map[string]interface{} `json:"external_details"`
}

func (r *UpdateKYCRequest) UnmarshalJSON(data []byte) error {
	type t UpdateKYCRequest
	var tt t
	if err := json.Unmarshal(data, &tt); err != nil {
		return err
	}
	*r = UpdateKYCRequest(tt)

	// marshal map back to json
	rawKYC, err := json.Marshal(r.KYCData)
	if err != nil {
		return err
	}

	// finally unmarshal to proper struct
	if err := json.Unmarshal(rawKYC, &r.KYCDataStruct); err != nil {
		return err
	}

	return nil
}

type AccountTypeToSet struct {
	Int    int    `json:"int"`
	String string `json:"string"`
}

type UpdateSaleDetailsRequest struct {
	SaleID     uint64                 `json:"sale_id"`
	NewDetails map[string]interface{} `json:"new_details"`
}

type UpdateSaleEndTimeRequest struct {
	SaleID     uint64 `json:"sale_id"`
	NewEndTime Time   `json:"new_end_time"`
}

type WithdrawalRequest struct {
	BalanceID              string                 `json:"balance_id"`
	Amount                 Amount                 `json:"amount"`
	FixedFee               Amount                 `json:"fixed_fee"`
	PercentFee             Amount                 `json:"percent_fee"`
	PreConfirmationDetails map[string]interface{} `json:"pre_confirmation_details"`
	ExternalDetails        map[string]interface{} `json:"external_details"`
	DestAssetCode          string                 `json:"dest_asset_code"`
	DestAssetAmount        Amount                 `json:"dest_asset_amount"`
	ReviewerDetails        map[string]interface{} `json:"reviewer_details"`
}
