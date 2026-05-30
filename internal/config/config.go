package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	GeoIPCityDB     string
	GeoIPASNDB      string
	AbuseIPDBKey    string
	RedisAddr       string
	RedisPassword   string
	RedisDB         int
	RateLimitIP     int
	RateLimitWHOIS  int
	RateLimitDomain int
	RateLimitEmail  int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using environment variables: %v", err)
	}

	return &Config{
		Port:            getEnv("PORT", "8080"),
		GeoIPCityDB:     getEnv("GEOIP_CITY_DB", "/app/data/GeoLite2-City.mmdb"),
		GeoIPASNDB:      getEnv("GEOIP_ASN_DB", "/app/data/GeoLite2-ASN.mmdb"),
		AbuseIPDBKey:    getEnv("ABUSEIPDB_KEY", ""),
		RedisAddr:       getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		RedisDB:         getEnvInt("REDIS_DB", 0),
		RateLimitIP:     getEnvInt("RATE_LIMIT_IP", 100),
		RateLimitWHOIS:  getEnvInt("RATE_LIMIT_WHOIS", 20),
		RateLimitDomain: getEnvInt("RATE_LIMIT_DOMAIN", 20),
		RateLimitEmail:  getEnvInt("RATE_LIMIT_EMAIL", 20),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		log.Printf("Invalid int for %s: %v, using default %d", key, err, defaultVal)
		return defaultVal
	}
	return n
}
