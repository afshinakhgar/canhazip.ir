package admin

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"iplocation.sabaai.ir/internal/infrastructure/blocklist"
	"iplocation.sabaai.ir/internal/infrastructure/requestlog"
)

//go:embed ui/index.html
var indexHTML []byte

// Handler holds dependencies for admin HTTP handlers.
type Handler struct {
	logger    *requestlog.Logger
	blocklist *blocklist.Checker
}

// NewHandler creates a Handler wired to the given logger and blocklist checker.
func NewHandler(logger *requestlog.Logger, bl *blocklist.Checker) *Handler {
	return &Handler{logger: logger, blocklist: bl}
}

// Register wires all admin routes onto the provided RouterGroup.
// The group should already have TokenAuth middleware applied.
func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.GET("/", h.serveIndex)
	rg.GET("/api/stats", h.getStats)
	rg.GET("/api/requests", h.getRequests)
	rg.POST("/api/blocklist/reload", h.reloadBlocklist)
	rg.POST("/api/blocklist/upload", h.uploadBlocklist)
}

// GET /admin/
func (h *Handler) serveIndex(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
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
