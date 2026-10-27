package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Host string `env:"RATE_LIMITER_HOST" envDefault:"127.0.0.1"`
	Port string `env:"RATE_LIMITER_PORT" envDefault:"8000"`
}

func New() (*Config, error) {

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("env.ParseAs: %w", err)
	}

	return &cfg, nil

}
