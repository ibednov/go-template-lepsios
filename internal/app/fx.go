package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ibednov/go-lepsios/log"
	"github.com/ibednov/go-template-lepsios/internal/config"
	"github.com/ibednov/go-template-lepsios/internal/domain/example"
	sidecarfiles "github.com/ibednov/go-template-lepsios/internal/sidecar/files"
	sidecarpostgres "github.com/ibednov/go-template-lepsios/internal/sidecar/postgres"
	sidecarredis "github.com/ibednov/go-template-lepsios/internal/sidecar/redis"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var Module = fx.Options(
	config.Module,
	sidecarfiles.Module,
	sidecarpostgres.Module,
	sidecarredis.Module,
	example.Module,
	fx.Provide(
		NewLogger,
		NewEngine,
		NewHTTPServer,
		fx.Annotate(
			NewReadinessChecker,
			fx.ResultTags(`name:"readiness"`),
		),
	),
	fx.Invoke(
		RegisterRoutes,
		StartServer,
	),
)

func Run() {
	fx.New(Module).Run()
}

func NewLogger(cfg *config.Config) (zerolog.Logger, error) {
	return log.Setup(log.Config{
		Env:         cfg.Environment,
		ServiceName: cfg.ServiceName,
		Level:       cfg.LogLevel,
		VectorURL:   cfg.LogVectorURL,
	})
}

func NewHTTPServer(cfg *config.Config, engine *gin.Engine) *http.Server {
	return &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
}

func StartServer(lc fx.Lifecycle, srv *http.Server, logger zerolog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info().Str("addr", srv.Addr).Msg("http server starting")
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal().Err(err).Msg("http server failed")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			logger.Info().Msg("http server stopping")
			if err := srv.Shutdown(shutdownCtx); err != nil {
				return fmt.Errorf("shutdown: %w", err)
			}
			return nil
		},
	})
}
