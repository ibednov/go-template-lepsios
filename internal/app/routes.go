package app

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	lepsioshttpx "github.com/ibednov/go-lepsios/httpx"
	examplehandler "github.com/ibednov/go-template-lepsios/internal/domain/example/transport/api/http/v1/handler"
	"github.com/ibednov/go-template-lepsios/internal/sidecar/postgres"
	"github.com/ibednov/go-template-lepsios/internal/sidecar/redis"
	"go.uber.org/fx"
)

type routeParams struct {
	fx.In

	Engine    *gin.Engine
	Postgres  *postgres.Pool
	Redis     *redis.Client
	Readiness lepsioshttpx.Checker `name:"readiness"`
	Example   *examplehandler.Handler
}

func RegisterRoutes(p routeParams) {
	lepsioshttpx.RegisterHealth(
		p.Engine,
		lepsioshttpx.NewChecker("liveness", func(context.Context) error { return nil }),
		p.Readiness,
	)
	lepsioshttpx.RegisterMetrics(p.Engine)

	v1 := p.Engine.Group("/api/v1")
	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	p.Example.Register(v1.Group("/example"))
}

func NewReadinessChecker(pg *postgres.Pool, rdb *redis.Client) lepsioshttpx.Checker {
	return lepsioshttpx.NewChecker("deps", func(ctx context.Context) error {
		if pg.Ready() {
			if err := pg.Ping(ctx); err != nil {
				return fmt.Errorf("postgres: %w", err)
			}
		}
		if rdb.Ready() {
			if err := rdb.Ping(ctx); err != nil {
				return fmt.Errorf("redis: %w", err)
			}
		}
		return nil
	})
}
