package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// StoredAuth is persisted in the config file
type StoredAuth struct {
	AccessToken     string    `json:"access_token"`
	RefreshToken    string    `json:"refresh_token"`
	ExpiresAt       time.Time `json:"expires_at"`
	UserUUID        string    `json:"user_uuid"`
	DefaultListUUID string    `json:"default_list_uuid"`
	DefaultListName string    `json:"default_list_name"`
	Email           string    `json:"email"`
	Name            string    `json:"name"`
}

func Dir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "bring-tui")
}

func Path() string {
	return filepath.Join(Dir(), "config.json")
}

func Load() (*StoredAuth, error) {
	data, err := os.ReadFile(Path())
	if err != nil {
		return nil, err
	}
	var auth StoredAuth
	if err := json.Unmarshal(data, &auth); err != nil {
		return nil, err
	}
	return &auth, nil
}

func Save(auth *StoredAuth) error {
	if err := os.MkdirAll(Dir(), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(auth, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(Path(), data, 0600)
}

func Delete() error {
	return os.Remove(Path())
}
