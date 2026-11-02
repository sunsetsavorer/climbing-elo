package config

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

const DefaultEnvFile = "configs/.env"

type (
	Config struct {
		Database Database `envPrefix:"DATABASE_"`
		Bot      Bot      `envPrefix:"BOT_"`
	}

	Database struct {
		DSN string `env:"DSN,required,notEmpty"`
	}

	Bot struct {
		Token string `env:"TOKEN,required,notEmpty"`
	}
)

// Load reads configuration from environment variables.
// If envFile exists, its values are loaded into the environment first;
// variables already set in the environment take precedence.
func Load(envFile string) (*Config, error) {
	if err := godotenv.Load(envFile); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("load env file %q: %w", envFile, err)
	}

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &cfg, nil
}
