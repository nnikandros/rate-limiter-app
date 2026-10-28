package clientdb

type Repository interface {
	GetClientConfig(clientID string) (ClientConfig, bool)
}
