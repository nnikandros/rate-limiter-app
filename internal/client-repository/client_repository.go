package clientrepository

import (
	"log"
	"rate-limiter-app/internal/config"
)

type ClientRepository interface {
	ListClients() []string
}

func New(cfg config.Config) ClientRepository {
	if cfg.UserRepository != "environment" {
		log.Fatal("not implemented")
	}

	return envUsersRepository{}

}
