package middleware

import (
	"net"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"iplocation.sabaai.ir/internal/infrastructure/requestlog"
)

// RequestLogger returns a Gin middleware that records each request into the
// provided requestlog.Logger.
func RequestLogger(logger *requestlog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start).Milliseconds()

		ip := extractRequestIP(c)
		entry := requestlog.Entry{
			Timestamp: start.UTC().Format(time.RFC3339),
			IP:        ip,
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Status:    c.Writer.Status(),
			LatencyMs: latency,
			UserAgent: c.Request.UserAgent(),
		}
		logger.Push(entry)
	}
}

// extractRequestIP returns the best available client IP for logging.
func extractRequestIP(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if ip := strings.TrimSpace(parts[0]); ip != "" {
			return ip
		}
	}
	host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return host
}
