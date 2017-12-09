package create_account_streamer

type OperationType int

const (
	OperationsTypeCreateAccount  OperationType = 0
	OperationsTypeForfeitRequest OperationType = 8
)

var (
	operationsTypeNames = map[OperationType]string{
		OperationsTypeCreateAccount:  "create_account",
		OperationsTypeForfeitRequest: "forfeit_request",
	}
)

// String returns string representation of OperationType if this type is unknown.
func (ot OperationType) String() string {
	name, ok := operationsTypeNames[ot]
	if ok {
		return name
	}

	return "unknown_operation_type"
}
