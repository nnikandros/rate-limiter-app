package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Host           string `env:"RATE_LIMITER_HOST" envDefault:"127.0.0.1"`
	Port           string `env:"RATE_LIMITER_PORT" envDefault:"8000"`
	UserRepository string `env:"RATE_LIMITER_USER_REPO" envDefault:"environment"`
	// UserDbUrl      string `env:"RATE_LIMITER_USER_DB_URL"`
	SqliteDbPath string `env:"RATE_LIMITER_SQLITE_DB" envDefault:"ratelimiter.db"`
	MockDbPath   string `env:"RATE_LIMITER_MOCK_DB" envDefault:"mock_db.json"`
}

func New() (*Config, error) {

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("env.ParseAs: %w", err)
	}

	return &cfg, nil

}
