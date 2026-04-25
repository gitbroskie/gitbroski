// Package mr provides MR/PR management functionality.
package mr

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"gitbroski/utils/logger"
)

// MR represents a saved merge request/pull request.
type MR struct {
	URL     string    `json:"url"`
	Title   string    `json:"title,omitempty"`
	Service string    `json:"service,omitempty"`
	SavedAt time.Time `json:"saved_at"`
}

// Store holds all saved MRs.
type Store struct {
	MRs []MR `json:"mrs"`
}

func getStorePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Error("Could not get home directory: " + err.Error())
		return ""
	}
	return filepath.Join(homeDir, ".config", "gitbroski", "mrs.json")
}

// LoadStore loads the MR store from disk.
func LoadStore() (*Store, error) {
	storePath := getStorePath()
	if storePath == "" {
		return &Store{MRs: []MR{}}, nil
	}

	data, err := os.ReadFile(storePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Store{MRs: []MR{}}, nil
		}
		return nil, err
	}

	var store Store
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	return &store, nil
}

// SaveStore saves the MR store to disk.
func SaveStore(store *Store) error {
	storePath := getStorePath()
	if storePath == "" {
		return nil
	}

	// Ensure directory exists
	dir := filepath.Dir(storePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(storePath, data, 0o600)
}
