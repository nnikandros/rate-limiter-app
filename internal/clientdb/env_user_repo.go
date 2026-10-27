package clientdb

import (
	"os"
	"rate-limiter-app/internal/clientdb/queries"
)

var users = os.Getenv("RATE_LIMITER_USERS")

type envUsersRepository struct{}

func (e envUsersRepository) ListClients() []queries.Client {

	return nil

}
