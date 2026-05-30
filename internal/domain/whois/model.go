package whois

// WhoisInfo holds parsed WHOIS data for a domain.
type WhoisInfo struct {
	Domain      string   `json:"domain"`
	Registrar   string   `json:"registrar"`
	Created     string   `json:"created"`
	Updated     string   `json:"updated"`
	Expires     string   `json:"expires"`
	Status      []string `json:"status"`
	Nameservers []string `json:"nameservers"`
	RawText     string   `json:"raw_text"`
}

// Repository is the port for WHOIS lookups.
type Repository interface {
	Lookup(domain string) (*WhoisInfo, error)
}
