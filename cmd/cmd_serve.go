package cmd

import (
	"log"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/server"
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

		log.Printf("Listening on address %s", cfg.App().Address)

		return server.InitServer(cfg)
	},
}
