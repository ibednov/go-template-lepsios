package config

import (
	"log"
	"strings"
	"time"

	lepsiosconfig "github.com/ibednov/go-lepsios/config"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(Load),
)

type Config struct {
	Port               string   `envconfig:"PORT" default:"8080"`
	Environment        string   `envconfig:"ENVIRONMENT" default:"development"`
	LogLevel           string   `envconfig:"LOG_LEVEL" default:"info"`
	LogVectorURL       string   `envconfig:"LOG_VECTOR_URL" default:""`
	ServiceName        string   `envconfig:"SERVICE_NAME" default:"go-template-lepsios"`
	CORSAllowedOrigins []string `envconfig:"CORS_ALLOWED_ORIGINS" default:"*"`
	Server             ServerConfig
	Database           DatabaseConfig
	Redis              RedisConfig
	Files              FilesConfig
	Migration          MigrationConfig
}

type ServerConfig struct {
	ReadTimeout  time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"10s"`
	WriteTimeout time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"10s"`
	IdleTimeout  time.Duration `envconfig:"SERVER_IDLE_TIMEOUT" default:"60s"`
}

type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     string `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER" default:"template"`
	Password string `envconfig:"DB_PASSWORD" default:"template"`
	Name     string `envconfig:"DB_NAME" default:"template"`
	SSLMode  string `envconfig:"DB_SSL_MODE" default:"disable"`
	Optional bool   `envconfig:"DB_OPTIONAL" default:"false"`
}

type RedisConfig struct {
	Host     string `envconfig:"REDIS_HOST" default:"localhost"`
	Port     string `envconfig:"REDIS_PORT" default:"6379"`
	Password string `envconfig:"REDIS_PASSWORD" default:""`
	DB       int    `envconfig:"REDIS_DB" default:"0"`
	Optional bool   `envconfig:"REDIS_OPTIONAL" default:"false"`
}

// FilesConfig selects go-lepsios files.Adapter (local | s3 | minio).
type FilesConfig struct {
	Driver    string `envconfig:"FILES_DRIVER" default:"local"`
	LocalDir  string `envconfig:"FILES_LOCAL_DIR" default:"var/files"`
	Endpoint  string `envconfig:"FILES_S3_ENDPOINT" default:"localhost:9000"`
	AccessKey string `envconfig:"FILES_S3_ACCESS_KEY" default:"minioadmin"`
	SecretKey string `envconfig:"FILES_S3_SECRET_KEY" default:"minioadmin"`
	Bucket    string `envconfig:"FILES_S3_BUCKET" default:"template"`
	Region    string `envconfig:"FILES_S3_REGION" default:"us-east-1"`
	UseSSL    bool   `envconfig:"FILES_S3_USE_SSL" default:"false"`
	PublicURL string `envconfig:"FILES_S3_PUBLIC_URL" default:""`
}

// MigrationConfig controls entrypoint migrate + pre-migrate backup via files.Adapter.
type MigrationConfig struct {
	Run          bool   `envconfig:"MIGRATION_RUN" default:"true"`
	Backup       bool   `envconfig:"MIGRATION_BACKUP" default:"true"`
	BackupPrefix string `envconfig:"MIGRATION_BACKUP_PREFIX" default:"backups/pre-migrate"`
}

func (c DatabaseConfig) DSN() string {
	return "postgres://" + c.User + ":" + c.Password + "@" + c.Host + ":" + c.Port + "/" + c.Name + "?sslmode=" + c.SSLMode
}

func (c RedisConfig) Addr() string {
	return c.Host + ":" + c.Port
}

func Load() *Config {
	var cfg Config
	if err := lepsiosconfig.Load(&cfg); err != nil {
		log.Panicf("failed to load config: %v", err)
	}
	for i, origin := range cfg.CORSAllowedOrigins {
		cfg.CORSAllowedOrigins[i] = strings.TrimSpace(origin)
	}
	cfg.Files.Driver = strings.ToLower(strings.TrimSpace(cfg.Files.Driver))
	return &cfg
}
