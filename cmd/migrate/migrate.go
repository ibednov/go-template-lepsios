package migrate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	lepsiosmigrate "github.com/ibednov/go-lepsios/migrate"
	"github.com/ibednov/go-template-lepsios/internal/config"
	sidecarfiles "github.com/ibednov/go-template-lepsios/internal/sidecar/files"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Database migrations (Goose)",
	}
	cmd.AddCommand(upCmd())
	cmd.AddCommand(downCmd())
	return cmd
}

func upCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Apply all pending migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run("up")
		},
	}
}

func downCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Rollback the last migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run("down")
		},
	}
}

func run(direction string) error {
	cfg := config.Load()
	dir, err := migrationsDir()
	if err != nil {
		return err
	}
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	db, err := sql.Open("postgres", cfg.Database.DSN())
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer func() { _ = db.Close() }()
	if err := db.Ping(); err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}

	switch direction {
	case "up":
		if err := backupBeforeUp(cfg, db, dir); err != nil {
			return err
		}
		if err := goose.Up(db, dir); err != nil {
			return fmt.Errorf("migrate up: %w", err)
		}
	case "down":
		if err := goose.Down(db, dir); err != nil {
			return fmt.Errorf("migrate down: %w", err)
		}
	default:
		return fmt.Errorf("unknown direction %q", direction)
	}
	fmt.Printf("migrations %s OK (%s)\n", direction, dir)
	return nil
}

func backupBeforeUp(cfg *config.Config, db *sql.DB, migrationsDir string) error {
	current, err := goose.GetDBVersion(db)
	if err != nil {
		return fmt.Errorf("read current migration version: %w", err)
	}
	pending, err := goose.CollectMigrations(migrationsDir, current, goose.MaxVersion)
	if err != nil {
		return fmt.Errorf("collect pending migrations: %w", err)
	}
	if len(pending) == 0 {
		fmt.Println("migration backup: skipped (no pending migrations)")
		return nil
	}
	if !cfg.Migration.Backup {
		fmt.Printf("migration backup: DISABLED (MIGRATION_BACKUP=false), pending=%d — continuing without backup\n", len(pending))
		return nil
	}

	adapter, err := sidecarfiles.NewAdapter(cfg)
	if err != nil {
		return fmt.Errorf("files adapter for migration backup: %w", err)
	}

	mapped := make([]lepsiosmigrate.PendingMigration, 0, len(pending))
	for _, migration := range pending {
		mapped = append(mapped, lepsiosmigrate.PendingMigration{
			Version: migration.Version,
			Source:  path.Base(migration.Source),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	if err := lepsiosmigrate.BackupBeforeMigrate(ctx, lepsiosmigrate.BackupInput{
		Adapter: adapter,
		DB: lepsiosmigrate.DatabaseDSN{
			Host:     cfg.Database.Host,
			Port:     cfg.Database.Port,
			User:     cfg.Database.User,
			Password: cfg.Database.Password,
			Name:     cfg.Database.Name,
			SSLMode:  cfg.Database.SSLMode,
		},
		Prefix:           cfg.Migration.BackupPrefix,
		StorageTypeLabel: sidecarfiles.DriverLabel(cfg),
		CurrentVersion:   current,
		Pending:          mapped,
	}); err != nil {
		return fmt.Errorf("pre-migration backup failed; migrations were not applied: %w", err)
	}
	return nil
}

func migrationsDir() (string, error) {
	candidates := []string{"/app/migrations", "migrations"}
	if wd, err := os.Getwd(); err == nil {
		candidates = append([]string{
			filepath.Join(wd, "migrations"),
			filepath.Join(wd, "..", "migrations"),
		}, candidates...)
	}
	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && st.IsDir() {
			return c, nil
		}
	}
	return "", errors.New("migrations directory not found")
}
