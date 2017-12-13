package notificator

//go:generate go-codegen
//go:generate gofmt -w main_generated.go

// for new payload instances add embed `payload` interface
// run gb generate
// and god bless golang type system

type Payload interface {
	legitPayload()
}

type EmailRequestPayload struct {
	payload Payload

	Destination string `json:"destination"`
	Subject     string `json:"subject"`
	Message     string `json:"message"`
}

type SMSRequestPayload struct {
	payload Payload

	Destination string `json:"destination"`
	Message     string `json:"message"`
}
