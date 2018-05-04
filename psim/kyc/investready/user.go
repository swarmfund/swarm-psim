package investready

type User struct {
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	DateOfBirth int64      `json:"dob"`
	Hash        string     `json:"hash"`
	Status      UserStatus `json:"status"`
}

func (u User) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"first_name":    u.FirstName,
		"last_name":     u.LastName,
		"email":         u.Email,
		"date_of_birth": u.DateOfBirth,
		"hash":          u.Hash,
		"status":        u.Status,
	}
}

type UserStatus struct {
	Message    StatusMessage `json:"message"`
	Accredited int           `json:"accredited"`
	// The response struct also contains fields:
	// - via
	// - certificate_url
	// - verification_type ("3rd Party", ... ?)
	// - expires_on
}

func (u UserStatus) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"message":    u.Message,
		"accredited": u.Accredited,
	}
}

type StatusMessage string

const (
	AccreditedStatusMessage StatusMessage = "Accredited"
	PendingStatusMessage    StatusMessage = "Pending"
	DeniedStatusMessage     StatusMessage = "Denied"
)
