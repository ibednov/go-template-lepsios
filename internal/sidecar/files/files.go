package files

import (
	"fmt"
	"strings"

	"github.com/ibednov/go-lepsios/files"
	"github.com/ibednov/go-lepsios/files/local"
	"github.com/ibednov/go-lepsios/files/s3minio"
	"github.com/ibednov/go-template-lepsios/internal/config"
	"go.uber.org/fx"
)

// Module always provides go-lepsios files.Adapter.
var Module = fx.Module("sidecar.files",
	fx.Provide(NewAdapter),
)

// NewAdapter builds local or S3/MinIO storage from FILES_* config.
func NewAdapter(cfg *config.Config) (files.Adapter, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.Files.Driver)) {
	case "s3", "minio":
		return s3minio.New(
			cfg.Files.Endpoint,
			cfg.Files.AccessKey,
			cfg.Files.SecretKey,
			cfg.Files.Bucket,
			cfg.Files.Region,
			cfg.Files.UseSSL,
			cfg.Files.PublicURL,
		)
	case "local", "":
		return local.New(cfg.Files.LocalDir)
	default:
		return nil, fmt.Errorf("unsupported FILES_DRIVER %q (want local|s3|minio)", cfg.Files.Driver)
	}
}

// DriverLabel is a short storage tag for logs / migrate backup.
func DriverLabel(cfg *config.Config) string {
	d := strings.ToLower(strings.TrimSpace(cfg.Files.Driver))
	if d == "" {
		return "local"
	}
	return d
}
