package email

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	emaildomain "iplocation.sabaai.ir/internal/domain/email"
)

var emailRegexp = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

var disposableDomains = map[string]bool{
	"mailinator.com":        true,
	"tempmail.com":          true,
	"guerrillamail.com":     true,
	"10minutemail.com":      true,
	"throwaway.email":       true,
	"yopmail.com":           true,
	"sharklasers.com":       true,
	"guerrillamailblock.com": true,
	"grr.la":                true,
	"guerrillamail.info":    true,
	"trashmail.com":         true,
	"dispostable.com":       true,
	"spamgourmet.com":       true,
	"maildrop.cc":           true,
	"fakeinbox.com":         true,
	"mailnull.com":          true,
	"spamfree24.org":        true,
	"discard.email":         true,
	"tempr.email":           true,
	"crazymailing.com":      true,
}

// Checker implements domain/email.Repository.
type Checker struct{}

// NewChecker returns a new email checker.
func NewChecker() *Checker {
	return &Checker{}
}

// Check validates the email address and returns reputation information.
func (c *Checker) Check(emailAddr string) (*emaildomain.EmailInfo, error) {
	info := &emaildomain.EmailInfo{
		Email: emailAddr,
	}

	// Syntax check.
	if !emailRegexp.MatchString(emailAddr) {
		info.Valid = false
		info.Reputation = "invalid"
		return info, nil
	}
	info.Valid = true

	// Extract domain.
	parts := strings.SplitN(emailAddr, "@", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("email: malformed address: %s", emailAddr)
	}
	domain := strings.ToLower(parts[1])
	info.Domain = domain

	// Disposable check.
	if disposableDomains[domain] {
		info.Disposable = true
		info.Reputation = "disposable"
		return info, nil
	}

	// MX record check.
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		info.MXValid = false
		info.Reputation = "invalid"
		return info, nil
	}
	info.MXValid = true
	info.Reputation = "valid"

	return info, nil
}
