package application

import (
	"fmt"

	emaildomain "iplocation.sabaai.ir/internal/domain/email"
)

// EmailService wraps the email repository.
type EmailService struct {
	repo emaildomain.Repository
}

// NewEmailService creates a new EmailService.
func NewEmailService(repo emaildomain.Repository) *EmailService {
	return &EmailService{repo: repo}
}

// Check validates and scores the given email address.
func (s *EmailService) Check(email string) (*emaildomain.EmailInfo, error) {
	info, err := s.repo.Check(email)
	if err != nil {
		return nil, fmt.Errorf("email service: %w", err)
	}
	return info, nil
}
