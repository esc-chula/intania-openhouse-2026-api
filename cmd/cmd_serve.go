package cmd

import (
	"context"
	"log"
	"time"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/migrations"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/server"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/database"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := getConfigFromCmd(cmd)
		if err != nil {
			return err
		}

		db := database.NewPostgresDB(cfg.Database())

		// migrate up here
		err = func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			goose.SetBaseFS(migrations.Migrations)
			if err := goose.SetDialect("postgres"); err != nil {
				return err
			}
			if err := goose.UpContext(ctx, db.DB, "."); err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			return err
		}

		log.Printf("Listening on address %s", cfg.App().Address)

		return server.InitServer(cfg, db)
	},
}
