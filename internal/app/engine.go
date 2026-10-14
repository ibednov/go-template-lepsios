package app

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ibednov/go-lepsios/httpx/middleware"
	"github.com/ibednov/go-template-lepsios/internal/config"
	"github.com/rs/zerolog"
)

var skipLogPaths = []string{
	"/health",
	"/health/live",
	"/health/ready",
	"/metrics",
}

func NewEngine(cfg *config.Config, logger zerolog.Logger) *gin.Engine {
	mode := gin.ReleaseMode
	switch cfg.Environment {
	case "development", "local", "dev", "test":
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	engine := gin.New()
	engine.Use(middleware.Recovery())
	engine.Use(middleware.TraceID())
	engine.Use(middleware.AccessLog(logger, skipLogPaths...))
	engine.Use(corsMiddleware(cfg.CORSAllowedOrigins))
	engine.Use(middleware.Metrics())
	return engine
}

func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowAll := len(allowedOrigins) == 0 || (len(allowedOrigins) == 1 && allowedOrigins[0] == "*")
		if origin != "" {
			if allowAll {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Credentials", "true")
			} else {
				for _, allowed := range allowedOrigins {
					if origin == strings.TrimSpace(allowed) {
						c.Header("Access-Control-Allow-Origin", origin)
						c.Header("Access-Control-Allow-Credentials", "true")
						break
					}
				}
			}
			c.Header("Vary", "Origin")
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}
