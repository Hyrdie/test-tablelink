package migrations

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

//go:embed *.up.sql
var migrationFiles embed.FS

type Migration struct {
	Version string
	UpSQL   string
}

func RunMigrations(ctx context.Context, db *sqlx.DB) error {
	// Create migrations table if it doesn't exist
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	var appliedMigrations []string
	err = db.SelectContext(ctx, &appliedMigrations, "SELECT version FROM migrations")
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Read migration files
	migrations := make(map[string]Migration)
	err = fs.WalkDir(migrationFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".up.sql") {
			return nil
		}

		version := strings.TrimSuffix(filepath.Base(path), ".up.sql")
		content, err := migrationFiles.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", path, err)
		}

		migrations[version] = Migration{
			Version: version,
			UpSQL:   string(content),
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	// Sort migrations by version
	var versions []string
	for version := range migrations {
		versions = append(versions, version)
	}
	sort.Strings(versions)

	// Apply pending migrations
	for _, version := range versions {
		// Skip if already applied
		alreadyApplied := false
		for _, applied := range appliedMigrations {
			if applied == version {
				alreadyApplied = true
				break
			}
		}
		if alreadyApplied {
			continue
		}

		// Start transaction
		tx, err := db.BeginTxx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to start transaction for migration %s: %w", version, err)
		}

		// Apply migration
		_, err = tx.ExecContext(ctx, migrations[version].UpSQL)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to apply migration %s: %w", version, err)
		}

		// Record migration
		_, err = tx.ExecContext(ctx, "INSERT INTO migrations (version) VALUES ($1)", version)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", version, err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", version, err)
		}
	}

	return nil
}
