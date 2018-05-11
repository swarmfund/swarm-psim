package horizon

type SubmitFailedError struct {
	Msg        string
	StatusCode int
	TXCode     string
	OpCodes    []string
}

func (e SubmitFailedError) Error() string {
	return e.Msg
}

func NewSubmitFailedError(msg string, statusCode int, txCode string, opCodes []string) SubmitFailedError {
	return SubmitFailedError{
		Msg:        msg,
		StatusCode: statusCode,
		TXCode:     txCode,
		OpCodes:    opCodes,
	}
}
