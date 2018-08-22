package issuance

const (
	TaskNotEnoughPreissuance    uint32 = 1
	TaskManualReviewAssetPolicy uint32 = 2
	TaskIssuanceLimitExceeded   uint32 = 4

	TaskVerifyDeposit uint32 = 1024
)

var (
	DefaultIssuanceTasks = TaskVerifyDeposit
)

// RequestOpt is the gathering structure for arguments to be set into CreateIssuanceRequest Operation.
// DEPRECATED: use xdrbuild.CreateIssuanceRequest directly
type RequestOpt struct {
	Reference string
	Receiver  string
	Asset     string
	Amount    uint64
	Details   string
	// AllTasks is the mask for Tasks to set on creation of the IssuanceRequest.
	AllTasks uint32
}

func (i RequestOpt) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"reference": i.Reference,
		"receiver":  i.Receiver,
		"asset":     i.Asset,
		"amount":    i.Amount,
		"details":   i.Details,
		"all_tasks": i.AllTasks,
	}
}
