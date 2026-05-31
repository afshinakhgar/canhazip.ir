package http

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"iplocation.sabaai.ir/internal/application"
	"iplocation.sabaai.ir/internal/config"
	"iplocation.sabaai.ir/internal/infrastructure/blocklist"
	"iplocation.sabaai.ir/internal/infrastructure/requestlog"
	"iplocation.sabaai.ir/internal/interfaces/http/admin"
	"iplocation.sabaai.ir/internal/interfaces/http/handlers"
	"iplocation.sabaai.ir/internal/interfaces/http/middleware"
)

// NewRouter builds and returns a configured Gin engine.
func NewRouter(
	cfg *config.Config,
	rdb *redis.Client,
	ipSvc *application.IPService,
	domainSvc *application.DomainService,
	whoisSvc *application.WhoisService,
	emailSvc *application.EmailService,
	reqLogger *requestlog.Logger,
	blChecker *blocklist.Checker,
	adminToken string,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())
	r.Use(middleware.RequestLogger(reqLogger))

	ipHandler := handlers.NewIPHandler(ipSvc)
	domainHandler := handlers.NewDomainHandler(domainSvc)
	whoisHandler := handlers.NewWhoisHandler(whoisSvc)
	emailHandler := handlers.NewEmailHandler(emailSvc)

	// GET / — caller IP info (rate limited per cfg.RateLimitIP).
	r.GET("/", middleware.RateLimit(rdb, cfg.RateLimitIP), ipHandler.GetCallerIP)

	// GET /:ip — specific IP info.
	// Note: Gin resolves static routes before parametric ones.
	r.GET("/:ip", middleware.RateLimit(rdb, cfg.RateLimitIP), ipHandler.GetIP)

	// GET /domain/:domain
	r.GET("/domain/:domain", middleware.RateLimit(rdb, cfg.RateLimitDomain), domainHandler.GetDomain)

	// GET /whois/:domain
	r.GET("/whois/:domain", middleware.RateLimit(rdb, cfg.RateLimitWHOIS), whoisHandler.GetWhois)

	// GET /email/:email
	r.GET("/email/:email", middleware.RateLimit(rdb, cfg.RateLimitEmail), emailHandler.GetEmail)

	// Admin panel — no rate limiting, protected by token auth.
	adm := r.Group("/admin")
	adm.Use(admin.TokenAuth(adminToken))
	admin.NewHandler(reqLogger, blChecker).Register(adm)

	return r
}

// corsMiddleware adds permissive CORS headers to every response.
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
