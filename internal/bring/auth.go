package bring

import (
	"fmt"
	"os"
	"time"

	"github.com/paulleonhardhellweg/bring-tui/internal/config"
)

// Authenticate resolves credentials in priority order:
// 1. Environment variables (BRING_EMAIL + BRING_PASSWORD)
// 2. Stored config with valid/refreshable token
// Returns (client, error). If no credentials found, returns ErrNeedsLogin.
func Authenticate() (*Client, *config.StoredAuth, error) {
	// 1. Try env vars
	email := os.Getenv("BRING_EMAIL")
	password := os.Getenv("BRING_PASSWORD")
	if email != "" && password != "" {
		return loginAndStore(email, password)
	}

	// 2. Try stored config
	stored, err := config.Load()
	if err == nil && stored.AccessToken != "" {
		// Check if token is still valid (with 5 min buffer)
		if time.Now().Before(stored.ExpiresAt.Add(-5 * time.Minute)) {
			client := NewClient(stored.AccessToken, stored.UserUUID)
			return client, stored, nil
		}

		// Try refresh
		if stored.RefreshToken != "" {
			tokenResp, err := RefreshToken(stored.RefreshToken, stored.AccessToken, stored.UserUUID)
			if err == nil {
				stored.AccessToken = tokenResp.AccessToken
				stored.RefreshToken = tokenResp.RefreshToken
				stored.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
				if err := config.Save(stored); err != nil {
					return nil, nil, fmt.Errorf("failed to save refreshed token: %w", err)
				}
				client := NewClient(stored.AccessToken, stored.UserUUID)
				return client, stored, nil
			}
		}
	}

	// 3. No credentials available
	return nil, nil, ErrNeedsLogin
}

// LoginAndStore authenticates with email/password and persists the token
func LoginAndStore(email, password string) (*Client, *config.StoredAuth, error) {
	return loginAndStore(email, password)
}

func loginAndStore(email, password string) (*Client, *config.StoredAuth, error) {
	auth, err := Login(email, password)
	if err != nil {
		return nil, nil, err
	}

	stored := &config.StoredAuth{
		AccessToken:     auth.AccessToken,
		RefreshToken:    auth.RefreshToken,
		ExpiresAt:       time.Now().Add(time.Duration(auth.ExpiresIn) * time.Second),
		UserUUID:        auth.UUID,
		DefaultListUUID: auth.BringListUUID,
		Email:           auth.Email,
		Name:            auth.Name,
	}

	if err := config.Save(stored); err != nil {
		return nil, nil, fmt.Errorf("failed to save auth: %w", err)
	}

	client := NewClient(auth.AccessToken, auth.UUID)
	return client, stored, nil
}

// ErrNeedsLogin indicates no stored credentials and no env vars
var ErrNeedsLogin = fmt.Errorf("no credentials found - please login")
