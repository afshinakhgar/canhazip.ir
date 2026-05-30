package whois

import (
	"context"
	"fmt"
	"time"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
	whoisdomain "iplocation.sabaai.ir/internal/domain/whois"
)

const timeout = 15 * time.Second

// Client implements domain/whois.Repository using likexian/whois.
type Client struct{}

// NewClient returns a new WHOIS client.
func NewClient() *Client {
	return &Client{}
}

// Lookup fetches and parses WHOIS data for the given domain.
func (c *Client) Lookup(domain string) (*whoisdomain.WhoisInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Run WHOIS query in a goroutine so we respect the context deadline.
	type result struct {
		raw string
		err error
	}
	ch := make(chan result, 1)
	go func() {
		raw, err := whois.Whois(domain)
		ch <- result{raw, err}
	}()

	var raw string
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("whois: timeout for domain %s", domain)
	case r := <-ch:
		if r.err != nil {
			return nil, fmt.Errorf("whois: lookup failed: %w", r.err)
		}
		raw = r.raw
	}

	parsed, err := whoisparser.Parse(raw)
	if err != nil {
		// Return raw-only result if parsing fails.
		return &whoisdomain.WhoisInfo{
			Domain:  domain,
			RawText: raw,
		}, nil
	}

	info := &whoisdomain.WhoisInfo{
		Domain:  domain,
		RawText: raw,
	}

	if parsed.Registrar != nil {
		info.Registrar = parsed.Registrar.Name
	}
	if parsed.Domain != nil {
		info.Created = parsed.Domain.CreatedDate
		info.Updated = parsed.Domain.UpdatedDate
		info.Expires = parsed.Domain.ExpirationDate
		info.Status = parsed.Domain.Status
		info.Nameservers = parsed.Domain.NameServers
	}

	return info, nil
}
