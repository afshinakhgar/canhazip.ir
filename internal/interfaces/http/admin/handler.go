package admin

import (
	_ "embed"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"iplocation.sabaai.ir/internal/infrastructure/banlist"
	"iplocation.sabaai.ir/internal/infrastructure/blocklist"
	"iplocation.sabaai.ir/internal/infrastructure/requestlog"
)

//go:embed ui/index.html
var indexHTML []byte

// Handler holds dependencies for admin HTTP handlers.
type Handler struct {
	logger    *requestlog.Logger
	blocklist *blocklist.Checker
	banlist   *banlist.BanList
}

// NewHandler creates a Handler wired to the given logger, blocklist checker, and ban list.
func NewHandler(logger *requestlog.Logger, bl *blocklist.Checker, bans *banlist.BanList) *Handler {
	return &Handler{logger: logger, blocklist: bl, banlist: bans}
}

// ServeIndex serves the admin HTML page (no auth required).
func (h *Handler) ServeIndex(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
}

// RegisterAPI wires the protected API routes onto rg (already has TokenAuth).
func (h *Handler) RegisterAPI(rg *gin.RouterGroup) {
	rg.GET("/stats", h.getStats)
	rg.GET("/requests", h.getRequests)
	rg.POST("/blocklist/reload", h.reloadBlocklist)
	rg.POST("/blocklist/upload", h.uploadBlocklist)
	rg.GET("/bans", h.getBans)
	rg.POST("/bans", h.banIP)
	rg.DELETE("/bans/:ip", h.unbanIP)
}

// GET /admin/api/stats
func (h *Handler) getStats(c *gin.Context) {
	totalReqs := int64(0)
	var blStats gin.H

	if h.logger != nil {
		totalReqs = h.logger.TotalCount()
	}

	if h.blocklist != nil {
		l1, l2, l1p, l2p := h.blocklist.Stats()
		blStats = gin.H{
			"level1_cidrs": l1,
			"level2_cidrs": l2,
			"level1_path":  l1p,
			"level2_path":  l2p,
		}
	} else {
		blStats = gin.H{
			"level1_cidrs": 0,
			"level2_cidrs": 0,
			"level1_path":  "",
			"level2_path":  "",
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total_requests": totalReqs,
		"blocklist":      blStats,
	})
}

// GET /admin/api/requests?limit=50
func (h *Handler) getRequests(c *gin.Context) {
	limit := int64(50)
	if lStr := c.Query("limit"); lStr != "" {
		if n, err := strconv.ParseInt(lStr, 10, 64); err == nil && n > 0 {
			if n > 200 {
				n = 200
			}
			limit = n
		}
	}

	if h.logger == nil {
		c.JSON(http.StatusOK, []requestlog.Entry{})
		return
	}
	entries := h.logger.Recent(limit)
	if entries == nil {
		entries = []requestlog.Entry{}
	}
	c.JSON(http.StatusOK, entries)
}

// POST /admin/api/blocklist/reload
func (h *Handler) reloadBlocklist(c *gin.Context) {
	if h.blocklist == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "blocklist not initialised"})
		return
	}
	if err := h.blocklist.Reload(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	l1, l2, l1p, l2p := h.blocklist.Stats()
	c.JSON(http.StatusOK, gin.H{
		"level1_cidrs": l1,
		"level2_cidrs": l2,
		"level1_path":  l1p,
		"level2_path":  l2p,
	})
}

// POST /admin/api/blocklist/upload?level=1|2
func (h *Handler) uploadBlocklist(c *gin.Context) {
	if h.blocklist == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "blocklist not initialised"})
		return
	}

	levelStr := c.Query("level")
	if levelStr != "1" && levelStr != "2" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query param level must be 1 or 2"})
		return
	}

	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("file field missing: %v", err)})
		return
	}

	f, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("open upload: %v", err)})
		return
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("read upload: %v", err)})
		return
	}

	var destPath string
	if levelStr == "1" {
		destPath = h.blocklist.Level1Path()
	} else {
		destPath = h.blocklist.Level2Path()
	}

	if destPath == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "blocklist path not configured"})
		return
	}

	if err := os.WriteFile(destPath, data, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("write file: %v", err)})
		return
	}

	if err := h.blocklist.Reload(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("reload after upload: %v", err)})
		return
	}

	l1, l2, l1p, l2p := h.blocklist.Stats()
	c.JSON(http.StatusOK, gin.H{
		"level1_cidrs": l1,
		"level2_cidrs": l2,
		"level1_path":  l1p,
		"level2_path":  l2p,
	})
}

// GET /admin/api/bans
func (h *Handler) getBans(c *gin.Context) {
	if h.banlist == nil {
		c.JSON(http.StatusOK, []string{})
		return
	}
	ips, err := h.banlist.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if ips == nil {
		ips = []string{}
	}
	c.JSON(http.StatusOK, ips)
}

// POST /admin/api/bans — body: {"ip":"1.2.3.4"}
func (h *Handler) banIP(c *gin.Context) {
	if h.banlist == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "banlist not initialised"})
		return
	}
	var body struct {
		IP string `json:"ip"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}
	if net.ParseIP(body.IP) == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid IP address"})
		return
	}
	if err := h.banlist.Ban(body.IP); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "ip": body.IP})
}

// DELETE /admin/api/bans/:ip
func (h *Handler) unbanIP(c *gin.Context) {
	if h.banlist == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "banlist not initialised"})
		return
	}
	ip := c.Param("ip")
	if err := h.banlist.Unban(ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
