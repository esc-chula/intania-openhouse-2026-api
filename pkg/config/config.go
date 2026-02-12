package config

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config interface {
	App() App
	Database() Database
	Firebase() Firebase

	String() string
}

type App struct {
	Address        string   `mapstructure:"address"         validate:"required"`
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	IsProduction   bool     `mapstructure:"is_production"`
}

type Database struct {
	DSN string `mapstructure:"dsn" validate:"required"`
}

type Firebase struct {
	ServiceAccountKeyFile string `mapstructure:"service_account_key_file"`
}

// -------------------------------------------------------------------------- //

type config struct {
	AppCfg      App      `mapstructure:"app"`
	DatabaseCfg Database `mapstructure:"database"`
	FirebaseCfg Firebase `mapstructure:"firebase"`
}

func (c *config) App() App           { return c.AppCfg }
func (c *config) Database() Database { return c.DatabaseCfg }
func (c *config) Firebase() Firebase { return c.FirebaseCfg }

func (c *config) String() string {
	jsonBytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Fatalf("Failed to convert config to JSON format: %v", err)
	}

	return string(jsonBytes)
}

//go:embed config_template.yaml
var configTemplate []byte

func InitConfig(envFile string) (Config, error) {
	// Load environments from file if not empty
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	baseViper := viper.New()
	baseViper.SetConfigType("yaml")
	baseViper.AutomaticEnv()
	baseViper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := baseViper.ReadConfig(bytes.NewReader(configTemplate)); err != nil {
		return nil, fmt.Errorf("Failed to read yaml config file: %w", err)
	}

	var cfg config
	if err := baseViper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal config: %w", err)
	}

	if err := validator.New().Struct(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
