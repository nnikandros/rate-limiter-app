package clientdb

import (
	"log"
	"rate-limiter-app/internal/clientdb/queries"
	"rate-limiter-app/internal/config"
)

type Repository interface {
	ListClients() []queries.Client
}

func New(cfg config.Config) Repository {
	if cfg.UserRepository != "environment" {
		log.Fatal("not implemented")
	}

	return envUsersRepository{}

}
