package handlers

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"iplocation.sabaai.ir/internal/application"
)

// IPHandler handles IP-related endpoints.
type IPHandler struct {
	service *application.IPService
}

// NewIPHandler creates a new IPHandler.
func NewIPHandler(service *application.IPService) *IPHandler {
	return &IPHandler{service: service}
}

// GetCallerIP handles GET / — returns info for the requester's own IP.
func (h *IPHandler) GetCallerIP(c *gin.Context) {
	ip := extractClientIP(c)
	h.lookupAndRespond(c, ip)
}

// GetIP handles GET /:ip — returns info for the specified IP.
func (h *IPHandler) GetIP(c *gin.Context) {
	ipStr := c.Param("ip")
	if net.ParseIP(ipStr) == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid IP address"})
		return
	}
	h.lookupAndRespond(c, ipStr)
}

func (h *IPHandler) lookupAndRespond(c *gin.Context, ip string) {
	info, err := h.service.Lookup(ip)
	if err != nil {
		errMsg := err.Error()
		// GeoIP returns "address not found" for private / reserved IPs.
		if isPrivateOrReservedError(errMsg) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "private or reserved IP address",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}
	c.JSON(http.StatusOK, info)
}

// extractClientIP returns the best available client IP address.
func extractClientIP(c *gin.Context) string {
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

func isPrivateOrReservedError(msg string) bool {
	lower := strings.ToLower(msg)
	return strings.Contains(lower, "address not found") ||
		strings.Contains(lower, "no such host") ||
		strings.Contains(lower, "addressnotfounderror")
}
