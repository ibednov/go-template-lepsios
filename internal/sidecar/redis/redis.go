package redis

import (
	"context"
	"fmt"
	"time"

	lepsiosredis "github.com/ibednov/go-lepsios/redis"
	"github.com/ibednov/go-template-lepsios/internal/config"
	goredis "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var Module = fx.Module("sidecar.redis",
	fx.Provide(NewClient),
)

type Client struct {
	RDB *goredis.Client
}

func (c *Client) Ready() bool {
	return c != nil && c.RDB != nil
}

func NewClient(cfg *config.Config, lc fx.Lifecycle, logger zerolog.Logger) (*Client, error) {
	rdb := lepsiosredis.NewClient(lepsiosredis.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		if cfg.Redis.Optional {
			logger.Warn().Err(err).Str("addr", cfg.Redis.Addr()).Msg("redis unavailable; continuing without client")
			return &Client{}, nil
		}
		return nil, fmt.Errorf("redis: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return rdb.Close()
		},
	})

	logger.Info().Str("addr", cfg.Redis.Addr()).Int("db", cfg.Redis.DB).Msg("redis client ready")
	return &Client{RDB: rdb}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	if c == nil || c.RDB == nil {
		return fmt.Errorf("redis: not connected")
	}
	return c.RDB.Ping(ctx).Err()
}
