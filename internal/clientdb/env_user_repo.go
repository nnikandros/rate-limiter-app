package clientdb

import (
	"os"
	"strings"
)

var users = os.Getenv("RATE_LIMITER_USERS")

type envUsersRepository struct{}

func (e envUsersRepository) ListClients() []string {

	return strings.Split(users, ",")

}
