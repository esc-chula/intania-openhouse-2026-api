package cmd

import (
	"context"
	"log"
	"time"

	"github.com/esc-chula/intania-openhouse-2026-api/pkg/database"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/seed"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed database with initial data",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := getConfigFromCmd(cmd)
		if err != nil {
			return err
		}
		db := database.NewPostgresDB(cfg.Database())

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		activities := seed.GetActivitySeedData()

		log.Printf("Seeding %d activities...", len(activities))
		_, err = db.NewInsert().Model(&activities).Exec(ctx)
		if err != nil {
			return err
		}

		booths := seed.GetBoothSeedData()

		log.Printf("Seeding %d booths...", len(booths))
		_, err = db.NewInsert().Model(&booths).Exec(ctx)
		if err != nil {
			return err
		}

		log.Println("Seeding completed successfully.")
		return nil
	},
}
