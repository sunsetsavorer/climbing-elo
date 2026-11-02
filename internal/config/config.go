package config

import (
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/url"
	"strconv"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

const DefaultEnvFile = "configs/.env.local"

type (
	Config struct {
		Database Database `envPrefix:"DB_"`
		Bot      Bot      `envPrefix:"BOT_"`
	}

	Database struct {
		Driver string `env:"DRIVER,required,notEmpty"`
		Host   string `env:"HOST,required,notEmpty"`
		Port   uint16 `env:"PORT,required,notEmpty"`
		User   string `env:"USER,required,notEmpty"`
		Pass   string `env:"PASS,required,notEmpty"`
		Name   string `env:"NAME,required,notEmpty"`
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

// DSN builds a connection URL, e.g. postgres://user:pass@host:5432/name?sslmode=disable.
func (d Database) DSN() string {
	u := url.URL{
		Scheme:   d.Driver,
		User:     url.UserPassword(d.User, d.Pass),
		Host:     net.JoinHostPort(d.Host, strconv.Itoa(int(d.Port))),
		Path:     d.Name,
		RawQuery: "sslmode=disable",
	}
	return u.String()
}
