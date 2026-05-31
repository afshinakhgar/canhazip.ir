package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"iplocation.sabaai.ir/internal/infrastructure/banlist"
)

// BanCheck returns a Gin middleware that blocks requests from banned IPs.
// Admin paths (/admin/*) are exempt so a mis-ban cannot lock out the admin UI.
// If bl is nil the middleware is a no-op (fail-open).
func BanCheck(bl *banlist.BanList) gin.HandlerFunc {
	return func(c *gin.Context) {
		if bl == nil {
			c.Next()
			return
		}

		// Skip ban check for admin paths.
		if strings.HasPrefix(c.Request.URL.Path, "/admin") {
			c.Next()
			return
		}

		ip := extractClientIP(c)
		if bl.IsBanned(ip) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}
