package operation

type ReviewableRequestType string

const (
	WithdrawalsReviewableRequestType ReviewableRequestType = "withdrawals"
	KYCReviewableRequestType         ReviewableRequestType = "update_kyc"
)
