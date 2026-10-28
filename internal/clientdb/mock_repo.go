package clientdb

import (
	"encoding/json"
	"fmt"
	"os"
)

// ClientConfig is the per-client rate-limit configuration as stored in
// mock_db.json. Store selects which rate-limiting strategy to apply
// ("inmem" or "redis").
type ClientConfig struct {
	ClientID string  `json:"client_id"`
	Rps      float64 `json:"rps"`
	Burst    int     `json:"burst"`
	Store    string  `json:"store"`
}

type MockRepository struct {
	byClient map[string]ClientConfig
}

// NewMockRepository loads client configuration from the JSON file at path.
func NewMockRepository(path string) (*MockRepository, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read mock db %q: %w", path, err)
	}

	var clients []ClientConfig
	if err := json.Unmarshal(data, &clients); err != nil {
		return nil, fmt.Errorf("parse mock db %q: %w", path, err)
	}

	byClient := make(map[string]ClientConfig, len(clients))
	for _, c := range clients {
		byClient[c.ClientID] = c
	}

	return &MockRepository{byClient: byClient}, nil
}

// GetClientConfig returns the config for clientID and whether it was found.
func (m *MockRepository) GetClientConfig(clientID string) (ClientConfig, bool) {
	c, ok := m.byClient[clientID]
	return c, ok
}
