package cmd

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/migrations"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/database"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate up|down|reset|create",
	Short: "Migrate database",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("expect an argument")
		}

		cfg, err := getConfigFromCmd(cmd)
		if err != nil {
			return err
		}
		db := database.NewPostgresDB(cfg.Database())

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		// Migrate the database
		goose.SetBaseFS(migrations.Migrations)
		if err := goose.SetDialect("postgres"); err != nil {
			return err
		}

		switch args[0] {
		case "up":
			err := goose.UpContext(ctx, db.DB, ".")
			return err
		case "down":
			err := goose.DownContext(ctx, db.DB, ".")
			return err
		case "reset":
			err := goose.ResetContext(ctx, db.DB, ".")
			if err != nil && strings.Contains(err.Error(), "failed to get status of migrations") {
				return nil
			}
			return err
		case "create":
			name := ""
			if len(args) >= 2 {
				name = args[1]
			}
			// Somehow goose.Create didn't use the BaseFS, so i need to manually locate the migrations folder
			err := goose.Create(db.DB, "./internal/migrations", name, "sql")
			return err
		default:
			return errors.New("invalid migrate argument.")
		}
	},
}
