package geoip

import (
	"fmt"
	"net"
	"time"

	"github.com/oschwald/geoip2-golang"
	ipdomain "iplocation.sabaai.ir/internal/domain/ip"
)

// Reader implements domain/ip.Repository using MaxMind GeoIP2 databases.
type Reader struct {
	cityDB *geoip2.Reader
	asnDB  *geoip2.Reader
	abuse  ipdomain.Repository // nested for reputation; may be nil
}

// NewReader opens both City and ASN databases. The abuseRepo is used for
// reputation lookups and may be nil (reputation will be empty).
func NewReader(cityPath, asnPath string, abuseRepo ipdomain.Repository) (*Reader, error) {
	city, err := geoip2.Open(cityPath)
	if err != nil {
		return nil, fmt.Errorf("open city db: %w", err)
	}
	asn, err := geoip2.Open(asnPath)
	if err != nil {
		city.Close()
		return nil, fmt.Errorf("open asn db: %w", err)
	}
	return &Reader{cityDB: city, asnDB: asn, abuse: abuseRepo}, nil
}

// Close releases database file handles.
func (r *Reader) Close() {
	r.cityDB.Close()
	r.asnDB.Close()
}

// Lookup returns full IPInfo for the given IP string.
func (r *Reader) Lookup(ipStr string) (*ipdomain.IPInfo, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	city, err := r.cityDB.City(ip)
	if err != nil {
		return nil, fmt.Errorf("city lookup: %w", err)
	}

	asnNum, isp, err := r.LookupASN(ipStr)
	if err != nil {
		// Non-fatal: ASN info is best-effort.
		asnNum = 0
		isp = ""
	}

	info := &ipdomain.IPInfo{
		IP:        ipStr,
		Continent: city.Continent.Names["en"],
		Country: ipdomain.Country{
			Name: city.Country.Names["en"],
			Code: city.Country.IsoCode,
		},
		City:     city.City.Names["en"],
		Lat:      city.Location.Latitude,
		Lon:      city.Location.Longitude,
		Timezone: city.Location.TimeZone,
		ASN:      asnNum,
		ISP:      isp,
	}

	if r.abuse != nil {
		type result struct{ rep *ipdomain.Reputation }
		ch := make(chan result, 1)
		go func() {
			rep, err := r.abuse.LookupReputation(ipStr)
			if err != nil || rep == nil {
				ch <- result{}
				return
			}
			ch <- result{rep}
		}()
		select {
		case res := <-ch:
			if res.rep != nil {
				info.Reputation = *res.rep
			}
		case <-time.After(600 * time.Millisecond):
		}
	}

	return info, nil
}

// LookupASN returns ASN number and organisation name for the given IP.
func (r *Reader) LookupASN(ipStr string) (uint, string, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0, "", fmt.Errorf("invalid IP address: %s", ipStr)
	}
	record, err := r.asnDB.ASN(ip)
	if err != nil {
		return 0, "", fmt.Errorf("asn lookup: %w", err)
	}
	return record.AutonomousSystemNumber, record.AutonomousSystemOrganization, nil
}

// LookupReputation delegates to the abuse repository if set.
func (r *Reader) LookupReputation(ipStr string) (*ipdomain.Reputation, error) {
	if r.abuse != nil {
		return r.abuse.LookupReputation(ipStr)
	}
	return &ipdomain.Reputation{Score: 0, Label: "clean"}, nil
}
