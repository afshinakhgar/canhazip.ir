package application

import (
	"fmt"

	ipdomain "iplocation.sabaai.ir/internal/domain/ip"
)

// IPService implements domain/ip.Service.
type IPService struct {
	repo ipdomain.Repository
}

// NewIPService creates a new IPService.
func NewIPService(repo ipdomain.Repository) *IPService {
	return &IPService{repo: repo}
}

// Lookup returns full IP info for the given IP address string.
func (s *IPService) Lookup(ip string) (*ipdomain.IPInfo, error) {
	info, err := s.repo.Lookup(ip)
	if err != nil {
		return nil, fmt.Errorf("ip service: %w", err)
	}
	return info, nil
}
