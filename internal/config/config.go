package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Host           string `env:"RATE_LIMITER_HOST" envDefault:"127.0.0.1"`
	Port           string `env:"RATE_LIMITER_PORT" envDefault:"8000"`
	UserRepository string `env:"RATE_LIMITER_USER_REPO" envDefault:"environment"`
	// Users          string `env:"RATE_LIMITER_USERS"`
	UserDbUrl string `env:"RATE_LIMITER_USER_DB_URL"`
}

func New() (*Config, error) {

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("env.ParseAs: %w", err)
	}

	return &cfg, nil

}
