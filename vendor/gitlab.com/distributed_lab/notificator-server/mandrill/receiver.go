package mandrill

type Receiver struct {
	Email string `json:"email"`
}

func NewReceiver(email string) *Receiver {
	return &Receiver{
		Email: email,
	}
}
