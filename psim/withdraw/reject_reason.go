package withdraw

const (
	// Here is the full list of RejectReasons, which Service can set into `reject_reason` of Request in case of validation error(s).
	RejectReasonInvalidAddress  RejectReason = "invalid_btc_address"
	RejectReasonTooLittleAmount RejectReason = "too_little_amount"
)

type RejectReason string
