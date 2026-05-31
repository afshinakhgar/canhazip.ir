package abuseipdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	ipdomain "iplocation.sabaai.ir/internal/domain/ip"
)

const (
	apiURL  = "https://api.abuseipdb.com/api/v2/check"
	timeout = 500 * time.Millisecond
)

// Client calls the AbuseIPDB v2 API.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new AbuseIPDB client with the given API key.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

type abuseResponse struct {
	Data struct {
		AbuseConfidenceScore int `json:"abuseConfidenceScore"`
	} `json:"data"`
}

// Lookup is a stub that satisfies the Repository interface via delegation from
// geoip.Reader; the GeoIP reader itself does not call this.
func (c *Client) Lookup(_ string) (*ipdomain.IPInfo, error) {
	return nil, fmt.Errorf("abuseipdb: Lookup not implemented — use LookupReputation")
}

// LookupASN is not implemented on AbuseIPDB.
func (c *Client) LookupASN(_ string) (uint, string, error) {
	return 0, "", fmt.Errorf("abuseipdb: LookupASN not implemented")
}

// LookupReputation queries AbuseIPDB for the given IP's abuse confidence score.
func (c *Client) LookupReputation(ipStr string) (*ipdomain.Reputation, error) {
	if c.apiKey == "" {
		return &ipdomain.Reputation{Score: 0, Label: "clean"}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("abuseipdb: build request: %w", err)
	}

	q := req.URL.Query()
	q.Set("ipAddress", ipStr)
	q.Set("maxAgeInDays", "90")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Key", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("abuseipdb: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("abuseipdb: unexpected status %d", resp.StatusCode)
	}

	var result abuseResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("abuseipdb: decode response: %w", err)
	}

	score := result.Data.AbuseConfidenceScore
	label := scoreToLabel(score)

	return &ipdomain.Reputation{Score: score, Label: label}, nil
}

func scoreToLabel(score int) string {
	switch {
	case score < 25:
		return "clean"
	case score < 75:
		return "suspicious"
	default:
		return "malicious"
	}
}
