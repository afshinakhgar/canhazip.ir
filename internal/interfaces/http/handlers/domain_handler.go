package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"iplocation.sabaai.ir/internal/application"
	"iplocation.sabaai.ir/internal/domain/domaininfo"
)

// DomainHandler handles domain-related endpoints.
type DomainHandler struct {
	service *application.DomainService
}

// NewDomainHandler creates a new DomainHandler.
func NewDomainHandler(service *application.DomainService) *DomainHandler {
	return &DomainHandler{service: service}
}

type whoisSummary struct {
	Domain      string   `json:"domain"`
	Registrar   string   `json:"registrar"`
	Created     string   `json:"created"`
	Updated     string   `json:"updated"`
	Expires     string   `json:"expires"`
	Status      []string `json:"status"`
	Nameservers []string `json:"nameservers"`
}

type domainInfoResponse struct {
	Domain    string          `json:"domain"`
	ARecords  []string        `json:"a_records"`
	MXRecords []domaininfo.MXRecord `json:"mx_records"`
	NSRecords []string        `json:"ns_records"`
	Whois     *whoisSummary   `json:"whois,omitempty"`
}

func domainResponse(info *domaininfo.DomainInfo) domainInfoResponse {
	r := domainInfoResponse{
		Domain:    info.Domain,
		ARecords:  info.ARecords,
		MXRecords: info.MXRecords,
		NSRecords: info.NSRecords,
	}
	if info.Whois != nil {
		r.Whois = &whoisSummary{
			Domain:      info.Whois.Domain,
			Registrar:   info.Whois.Registrar,
			Created:     info.Whois.Created,
			Updated:     info.Whois.Updated,
			Expires:     info.Whois.Expires,
			Status:      info.Whois.Status,
			Nameservers: info.Whois.Nameservers,
		}
	}
	return r
}

// GetDomain handles GET /domain/:domain.
func (h *DomainHandler) GetDomain(c *gin.Context) {
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
	c.JSON(http.StatusOK, domainResponse(info))
}
