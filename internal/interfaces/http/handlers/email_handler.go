package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"iplocation.sabaai.ir/internal/application"
)

// EmailHandler handles email reputation endpoints.
type EmailHandler struct {
	service *application.EmailService
}

// NewEmailHandler creates a new EmailHandler.
func NewEmailHandler(service *application.EmailService) *EmailHandler {
	return &EmailHandler{service: service}
}

// GetEmail handles GET /email/:email.
func (h *EmailHandler) GetEmail(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	info, err := h.service.Check(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}
