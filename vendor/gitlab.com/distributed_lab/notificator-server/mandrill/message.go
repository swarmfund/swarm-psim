package mandrill

type Message struct {
	HTML      string     `json:"html"`
	Subject   string     `json:"subject"`
	FromEmail string     `json:"from_email"`
	FromName  string     `json:"from_name"`
	To        []*Receiver `json:"to"`
}

func NewMessage(fromEmail, fromName, subject, htmlContent string, to []*Receiver) *Message {
	return &Message{
		HTML:      htmlContent,
		Subject:   subject,
		FromEmail: fromEmail,
		FromName:  fromName,
		To:        to,
	}
}
