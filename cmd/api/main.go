package main

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"iplocation.sabaai.ir/internal/application"
	"iplocation.sabaai.ir/internal/config"
	"iplocation.sabaai.ir/internal/infrastructure/blocklist"
	"iplocation.sabaai.ir/internal/infrastructure/email"
	"iplocation.sabaai.ir/internal/infrastructure/geoip"
	whoisinfra "iplocation.sabaai.ir/internal/infrastructure/whois"
	httpinterface "iplocation.sabaai.ir/internal/interfaces/http"
)

func main() {
	// Load configuration from .env / environment.
	cfg := config.Load()

	// Load local Firehol blocklists for offline reputation checking.
	blChecker, err := blocklist.NewChecker(cfg.BlocklistLevel1, cfg.BlocklistLevel2)
	if err != nil {
		log.Printf("Warning: blocklist load failed: %v — reputation will be unknown", err)
		blChecker = nil
	} else {
		log.Printf("Blocklist loaded: %s / %s", cfg.BlocklistLevel1, cfg.BlocklistLevel2)
	}

	// Initialise GeoIP readers.
	geoReader, err := geoip.NewReader(cfg.GeoIPCityDB, cfg.GeoIPASNDB, blChecker)
	if err != nil {
		log.Fatalf("Failed to open GeoIP databases: %v", err)
	}
	defer geoReader.Close()
	log.Printf("GeoIP databases loaded from %s and %s", cfg.GeoIPCityDB, cfg.GeoIPASNDB)

	// Initialise Redis client.
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis ping failed: %v — rate limiting will fail-open", err)
	} else {
		log.Printf("Redis connected at %s", cfg.RedisAddr)
	}

	// Wire up repositories and services.
	ipSvc := application.NewIPService(geoReader)
	whoisRepo := whoisinfra.NewClient()
	whoisSvc := application.NewWhoisService(whoisRepo)
	domainSvc := application.NewDomainService(whoisRepo)
	emailRepo := email.NewChecker()
	emailSvc := application.NewEmailService(emailRepo)

	// Build router and start server.
	router := httpinterface.NewRouter(cfg, rdb, ipSvc, domainSvc, whoisSvc, emailSvc)

	addr := ":" + cfg.Port
	log.Printf("Starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
