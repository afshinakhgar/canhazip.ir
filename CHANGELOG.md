# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-05-31

### Added
- Full Go rewrite with Gin, hexagonal architecture
- IP lookup: geo (City/Country/Continent), ASN/ISP, lat/lon/timezone, AbuseIPDB reputation
- Domain lookup: DNS A/MX/NS records + WHOIS summary
- WHOIS lookup: full registrar/dates/status/nameservers + raw text
- Email reputation: syntax validation, MX check, disposable domain detection
- Redis sliding-window rate limiting per endpoint (configurable via .env)
- Docker + docker-compose + Nginx reverse proxy deployment
