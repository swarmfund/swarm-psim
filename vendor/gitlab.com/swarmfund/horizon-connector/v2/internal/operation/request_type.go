package operation

type ReviewableRequestType string

const (
	WithdrawalsReviewableRequestType ReviewableRequestType = "withdrawals"
	// TODO When KYC is ready
	//KYCReviewableRequestType ReviewableRequestType = "kyc"
)
