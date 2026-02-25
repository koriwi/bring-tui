package bring

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/paulleonhardhellweg/bring-tui/internal/config"
)


// Authenticate resolves credentials in priority order:
// 1. Environment variables (BRING_EMAIL + BRING_PASSWORD)
// 2. Stored config, always refreshing the token on startup
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
		return newClientFromStored(stored), stored, nil
	}

	// 3. No credentials available
	return nil, nil, ErrNeedsLogin
}

// newClientFromStored creates a Client wired up to persist token refreshes.
func newClientFromStored(stored *config.StoredAuth) *Client {
	c := &Client{
		http:         &http.Client{Timeout: 10 * time.Second},
		accessToken:  stored.AccessToken,
		refreshToken: stored.RefreshToken,
		userUUID:     stored.UserUUID,
	}
	c.onRefresh = func(accessToken, refreshToken string, expiresAt time.Time) error {
		stored.AccessToken = accessToken
		stored.RefreshToken = refreshToken
		stored.ExpiresAt = expiresAt
		return config.Save(stored)
	}
	return c
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

	return newClientFromStored(stored), stored, nil
}

// RefreshStoredToken loads the stored auth, refreshes the access token, and saves it back.
func RefreshStoredToken() error {
	stored, err := config.Load()
	if err != nil {
		return fmt.Errorf("no stored credentials: %w", err)
	}
	if stored.RefreshToken == "" {
		return fmt.Errorf("no refresh token stored - please login first")
	}
	tokenResp, err := RefreshToken(stored.RefreshToken, stored.AccessToken, stored.UserUUID)
	if err != nil {
		return fmt.Errorf("refresh failed: %w", err)
	}
	stored.AccessToken = tokenResp.AccessToken
	stored.RefreshToken = tokenResp.RefreshToken
	stored.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	return config.Save(stored)
}

// ErrNeedsLogin indicates no stored credentials and no env vars
var ErrNeedsLogin = fmt.Errorf("no credentials found - please login")
