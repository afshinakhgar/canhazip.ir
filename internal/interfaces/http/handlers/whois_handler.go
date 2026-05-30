package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"iplocation.sabaai.ir/internal/application"
)

// WhoisHandler handles WHOIS endpoints.
type WhoisHandler struct {
	service *application.WhoisService
}

// NewWhoisHandler creates a new WhoisHandler.
func NewWhoisHandler(service *application.WhoisService) *WhoisHandler {
	return &WhoisHandler{service: service}
}

// GetWhois handles GET /whois/:domain.
func (h *WhoisHandler) GetWhois(c *gin.Context) {
	domain := c.Param("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "domain is required"})
		return
	}

	info, err := h.service.Lookup(domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}
