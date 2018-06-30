package request_monitor

// from gitlab.com/swarmfund/horizon/db2/history

type ReviewableRequestState int

const (
	// ReviewableRequestStatePending - request was just created or updated
	ReviewableRequestStatePending ReviewableRequestState = iota + 1
	// ReviewableRequestStateCanceled - was canceled by requestor
	ReviewableRequestStateCanceled
	// ReviewableRequestStateApproved - was approved by reviewer
	ReviewableRequestStateApproved
	// ReviewableRequestStateRejected - was rejected by reviewer, but still can be updated
	ReviewableRequestStateRejected
	// ReviewableRequestStatePermanentlyRejected - was rejected by reviewer, can't be updated
	ReviewableRequestStatePermanentlyRejected
)
