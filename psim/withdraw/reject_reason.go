package withdraw

const (
	// Here is the full list of RejectReasons, which Service can set into `reject_reason` of Request in case of validation error(s).
	RejectReasonMissingAddress    RejectReason = "missing_address"
	RejectReasonAddressNotAString RejectReason = "address_not_a_string"
	RejectReasonInvalidAddress    RejectReason = "invalid_btc_address"
	RejectReasonTooLittleAmount   RejectReason = "too_little_amount"
)

type RejectReason string
