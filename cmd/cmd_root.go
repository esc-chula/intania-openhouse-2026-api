package cmd

import (
	"context"
	"errors"
	"log"

	"github.com/esc-chula/intania-openhouse-2026-api/pkg/config"
	"github.com/spf13/cobra"
)

type configContextKey struct{}

var RootCmd = &cobra.Command{
	Use:   "intania-openhouse-2026-api",
	Short: "CLIs command to run the intania-openhouse-2026-api",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		envFile, err := cmd.Flags().GetString("env-file")
		if err != nil {
			return err
		}

		cfg, err := config.InitConfig(envFile)
		if err != nil {
			return err
		}
		setConfigToCmd(cmd, cfg)

		log.Printf("Config: %s", cfg.String())

		return nil
	},
}

func init() {
	RootCmd.PersistentFlags().String("env-file", "", "environment file")
	RootCmd.AddCommand(serveCmd, migrateCmd)
}

func setConfigToCmd(cmd *cobra.Command, cfg config.Config) {
	cmd.SetContext(context.WithValue(cmd.Context(), configContextKey{}, cfg))
}

func getConfigFromCmd(cmd *cobra.Command) (config.Config, error) {
	cfg, ok := cmd.Context().Value(configContextKey{}).(config.Config)
	if !ok {
		return nil, errors.New("config not found")
	}
	return cfg, nil
}
