package email

// EmailInfo holds email reputation and validation data.
type EmailInfo struct {
	Email      string `json:"email"`
	Valid      bool   `json:"valid"`
	Domain     string `json:"domain"`
	MXValid    bool   `json:"mx_valid"`
	Disposable bool   `json:"disposable"`
	Reputation string `json:"reputation"` // "valid", "invalid", "disposable"
}

// Repository is the port for email checks.
type Repository interface {
	Check(email string) (*EmailInfo, error)
}
