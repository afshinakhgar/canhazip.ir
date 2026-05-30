package middleware

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	redisinfra "iplocation.sabaai.ir/internal/infrastructure/redis"
)

// RateLimit returns a Gin middleware that enforces a sliding-window rate limit
// of limit requests per minute per client IP using Redis.
func RateLimit(rdb *redis.Client, limit int) gin.HandlerFunc {
	limiter := redisinfra.NewRateLimiter(rdb)
	window := time.Minute

	return func(c *gin.Context) {
		clientIP := extractClientIP(c)
		key := clientIP

		allowed, err := limiter.Allow(c.Request.Context(), key, limit, window)
		if err != nil {
			// On Redis failure, allow the request through (fail-open).
			c.Next()
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate limit exceeded",
				"retry_after": 60,
			})
			return
		}

		c.Next()
	}
}

// extractClientIP returns the best available client IP.
func extractClientIP(c *gin.Context) string {
	// Trust X-Forwarded-For when behind a reverse proxy.
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		// Take the first (leftmost) address which is the original client.
		parts := strings.Split(xff, ",")
		if ip := strings.TrimSpace(parts[0]); ip != "" {
			return ip
		}
	}

	// Fall back to RemoteAddr, stripping the port.
	host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return host
}
