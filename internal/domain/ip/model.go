package ip

// IPInfo holds all information about an IP address.
type IPInfo struct {
	IP         string     `json:"ip"`
	Continent  string     `json:"continent"`
	Country    Country    `json:"country"`
	City       string     `json:"city"`
	Lat        float64    `json:"lat"`
	Lon        float64    `json:"lon"`
	Timezone   string     `json:"timezone"`
	ASN        uint       `json:"asn"`
	ISP        string     `json:"isp"`
	Reputation Reputation `json:"reputation"`
}

// Country holds country name and ISO code.
type Country struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// Reputation holds abuse score and label.
type Reputation struct {
	Score int    `json:"score"`
	Label string `json:"label"`
}

// Repository is the port for IP lookups.
type Repository interface {
	Lookup(ip string) (*IPInfo, error)
	LookupASN(ip string) (asn uint, isp string, err error)
	LookupReputation(ip string) (*Reputation, error)
}

// Service is the application-level IP service port.
type Service interface {
	Lookup(ip string) (*IPInfo, error)
}
