package horizon

import "time"

type ForfeitRequest struct {
	ID             string    `json:"paging_token"`
	Exchange       string    `json:"exchange"`
	PaymentID      string    `json:"payment_id"`
	PaymentState   uint32    `json:"payment_state"`
	Accepted       *bool     `json:"accepted"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	FromEmail      string    `json:"from_email"`
	ToEmail        string    `json:"to_email"`
	RequestType    int32     `json:"request_type"`
	PaymentDetails struct {
		Amount      string `json:"amount"`
		UserDetails string `json:"user_details"`
	} `json:"payment_details"`
}
