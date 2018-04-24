package emails

type Task struct {
	Destination string
	Subject     string
	Message     string
}

func (t Task) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"destination": t.Destination,
		"subject":     t.Subject,
	}
}
