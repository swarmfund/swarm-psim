package verifier

type Request struct {
	Envelope string `json:"envelope"`
}

func (r Request) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"envelope": r.Envelope,
	}
}
