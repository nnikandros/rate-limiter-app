package clientdb

import (
	"os"
	"rate-limiter-app/internal/clientdb/queries"
)

var users = os.Getenv("RATE_LIMITER_USERS")

type envUsersRepository struct {
	q *queries.Queries
}

func (e envUsersRepository) ListClients() []queries.Client {

	return nil

}

func (e envUsersRepository) GetLimitByClient(clientId string) queries.GetRateLimitByClientRow {

	return queries.GetRateLimitByClientRow{}
}
