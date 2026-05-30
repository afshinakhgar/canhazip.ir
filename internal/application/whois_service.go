package application

import (
	"fmt"

	whoisdomain "iplocation.sabaai.ir/internal/domain/whois"
)

// WhoisService wraps the WHOIS repository.
type WhoisService struct {
	repo whoisdomain.Repository
}

// NewWhoisService creates a new WhoisService.
func NewWhoisService(repo whoisdomain.Repository) *WhoisService {
	return &WhoisService{repo: repo}
}

// Lookup fetches WHOIS information for the given domain.
func (s *WhoisService) Lookup(domain string) (*whoisdomain.WhoisInfo, error) {
	info, err := s.repo.Lookup(domain)
	if err != nil {
		return nil, fmt.Errorf("whois service: %w", err)
	}
	return info, nil
}
