# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Phase 1: Full Go rewrite with Gin, hexagonal architecture
- IP lookup (geo + ASN + AbuseIPDB reputation)
- Domain lookup (DNS A/MX/NS + WHOIS)
- WHOIS lookup
- Email reputation (syntax + MX + disposable check)
- Redis-based per-IP rate limiting (sliding window)
- Docker + docker-compose + Nginx deployment
