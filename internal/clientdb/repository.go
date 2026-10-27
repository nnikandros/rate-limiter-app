package clientdb

import (
	"database/sql"
	"log"
	"rate-limiter-app/internal/clientdb/queries"
	"rate-limiter-app/internal/config"
)

type Repository interface {
	ListClients() []queries.Client
	GetLimitByClient(clientId string) queries.GetRateLimitByClientRow
}

func New(cfg config.Config) Repository {

	db, err := sql.Open("sqlite3", cfg.UserDbUrl)
	if err != nil {
		log.Fatal(err)
	}

	// queries.New(db)

	return envUsersRepository{q: queries.New(db)}

}
