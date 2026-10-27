package clients

import (
	"log"
	"rate-limiter-app/internal/config"
)

type Repository interface {
	ListClients() []string
}

func New(cfg config.Config) Repository {
	if cfg.UserRepository != "environment" {
		log.Fatal("not implemented")
	}

	return envUsersRepository{}

}
