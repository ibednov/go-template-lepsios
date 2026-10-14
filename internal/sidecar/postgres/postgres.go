package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	lepsiosdb "github.com/ibednov/go-lepsios/db"
	"github.com/ibednov/go-template-lepsios/internal/config"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var Module = fx.Module("sidecar.postgres",
	fx.Provide(NewPool),
)

type Pool struct {
	DB *sql.DB
}

func (p *Pool) Ready() bool {
	return p != nil && p.DB != nil
}

func NewPool(cfg *config.Config, lc fx.Lifecycle, logger zerolog.Logger) (*Pool, error) {
	db, err := lepsiosdb.ProvideDB(lepsiosdb.DBConfig{
		URL:             cfg.Database.DSN(),
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 30 * time.Minute,
	})
	if err != nil {
		if cfg.Database.Optional {
			logger.Warn().Err(err).Msg("postgres unavailable; continuing without pool")
			return &Pool{}, nil
		}
		return nil, fmt.Errorf("postgres: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return db.Close()
		},
	})

	logger.Info().
		Str("host", cfg.Database.Host).
		Str("db", cfg.Database.Name).
		Msg("postgres pool ready")
	return &Pool{DB: db}, nil
}

func (p *Pool) Ping(ctx context.Context) error {
	if p == nil || p.DB == nil {
		return fmt.Errorf("postgres: not connected")
	}
	return p.DB.PingContext(ctx)
}
