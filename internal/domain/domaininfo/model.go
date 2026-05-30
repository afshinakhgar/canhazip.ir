package domaininfo

import "iplocation.sabaai.ir/internal/domain/whois"

// DomainInfo holds DNS and WHOIS summary for a domain.
type DomainInfo struct {
	Domain    string      `json:"domain"`
	ARecords  []string    `json:"a_records"`
	MXRecords []MXRecord  `json:"mx_records"`
	NSRecords []string    `json:"ns_records"`
	Whois     *whois.WhoisInfo `json:"whois,omitempty"`
}

// MXRecord holds an MX entry.
type MXRecord struct {
	Host string `json:"host"`
	Pref uint16 `json:"pref"`
}

// Repository is the port for domain info lookups.
type Repository interface {
	Lookup(domain string) (*DomainInfo, error)
}
