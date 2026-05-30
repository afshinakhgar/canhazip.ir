# ORDERS

## Phase 1 — Full API Rewrite in Go

**Status:** In Progress  
**Date:** 2026-05-30

### Scope
Full rewrite from PHP prototype to Go. API-only service.

### Endpoints
| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Caller's IP full info |
| GET | `/:ip` | Specific IP info |
| GET | `/domain/:domain` | Domain DNS + WHOIS |
| GET | `/whois/:domain` | WHOIS data only |
| GET | `/email/:email` | Email reputation |

### IP Response Fields
- ip, continent, country (name, code), city, lat, lon, timezone
- ASN, ISP (from GeoLite2-ASN.mmdb)
- reputation (AbuseIPDB: score + label)

### Domain Response Fields
- domain, A records, MX records, NS records, WHOIS summary

### WHOIS Response Fields
- domain, registrar, created, updated, expires, status, nameservers

### Email Response Fields
- email, valid (syntax), domain, mx_valid, disposable, reputation

### Stack
- Go 1.22+
- Gin web framework
- Hexagonal Architecture (ports & adapters)
- Redis (rate limiting, sliding window)
- MaxMind GeoLite2 (City + ASN + Country)
- AbuseIPDB API

### Rate Limits (per-IP, sliding window)
- IP endpoints: 100 req/min (RATE_LIMIT_IP)
- WHOIS/domain/email: 20 req/min (RATE_LIMIT_WHOIS, RATE_LIMIT_DOMAIN, RATE_LIMIT_EMAIL)

### Infrastructure
- Docker + docker-compose
- Nginx reverse proxy
- Ubuntu VPS target

### Project Structure
```
cmd/api/main.go
internal/
  config/
  domain/ip|whois|domain|email/   (entities + port interfaces)
  application/                     (use cases)
  infrastructure/geoip|abuseipdb|whois|email|redis/
  interfaces/http/handlers|middleware/
nginx/nginx.conf
Dockerfile
docker-compose.yml
.env.example
CHANGELOG.md
ORDERS.md
```
