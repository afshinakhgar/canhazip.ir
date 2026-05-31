package admin

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// TokenAuth returns a middleware that checks a Bearer token in the Authorization
// header or a ?token= query parameter. Returns 401 JSON on failure.
// If token is empty, auth is skipped (dev mode).
func TokenAuth(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if token == "" {
			c.Next()
			return
		}

		// Check Authorization: Bearer <token>
		if auth := c.GetHeader("Authorization"); strings.HasPrefix(auth, "Bearer ") {
			if strings.TrimPrefix(auth, "Bearer ") == token {
				c.Next()
				return
			}
		}

		// Check cookie (preferred — survives CDN header stripping)
		if ck, err := c.Cookie("admin_token"); err == nil && ck == token {
			c.Next()
			return
		}

		// Check ?token= query param
		if qt := c.Query("token"); qt == token {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
	}
}
