package clientrepository

import (
	"os"
	"strings"
)

var users = os.Getenv("USERS")

type envUsersRepository struct{}

func (e envUsersRepository) ListClients() []string {

	return strings.Split(users, ",")

}
