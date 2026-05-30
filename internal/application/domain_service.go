package application

import (
	"fmt"
	"net"

	domaindomain "iplocation.sabaai.ir/internal/domain/domaininfo"
	whoisdomain "iplocation.sabaai.ir/internal/domain/whois"
)

// DomainService resolves DNS records and embeds a WHOIS summary.
type DomainService struct {
	whoisRepo whoisdomain.Repository
}

// NewDomainService creates a new DomainService.
func NewDomainService(whoisRepo whoisdomain.Repository) *DomainService {
	return &DomainService{whoisRepo: whoisRepo}
}

// Lookup returns DNS (A, MX, NS) and WHOIS summary for the given domain.
func (s *DomainService) Lookup(domain string) (*domaindomain.DomainInfo, error) {
	info := &domaindomain.DomainInfo{
		Domain:    domain,
		ARecords:  []string{},
		MXRecords: []domaindomain.MXRecord{},
		NSRecords: []string{},
	}

	// A records.
	addrs, err := net.LookupHost(domain)
	if err == nil {
		info.ARecords = addrs
	}

	// MX records.
	mxs, err := net.LookupMX(domain)
	if err == nil {
		for _, mx := range mxs {
			info.MXRecords = append(info.MXRecords, domaindomain.MXRecord{
				Host: mx.Host,
				Pref: mx.Pref,
			})
		}
	}

	// NS records.
	nss, err := net.LookupNS(domain)
	if err == nil {
		for _, ns := range nss {
			info.NSRecords = append(info.NSRecords, ns.Host)
		}
	}

	// WHOIS summary (best-effort).
	if s.whoisRepo != nil {
		whoisInfo, err := s.whoisRepo.Lookup(domain)
		if err == nil {
			info.Whois = whoisInfo
		}
	}

	// Return an error only if all DNS lookups failed and we have no data.
	if len(info.ARecords) == 0 && len(info.MXRecords) == 0 && len(info.NSRecords) == 0 {
		return nil, fmt.Errorf("domain service: no DNS records found for %s", domain)
	}

	return info, nil
}
