package bring

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is the Bring! API client
type Client struct {
	http        *http.Client
	accessToken string
	userUUID    string
}

// NewClient creates a new Bring! API client with auth credentials
func NewClient(accessToken, userUUID string) *Client {
	return &Client{
		http:        &http.Client{Timeout: 10 * time.Second},
		accessToken: accessToken,
		userUUID:    userUUID,
	}
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("X-BRING-API-KEY", APIKey)
	req.Header.Set("X-BRING-CLIENT", "webApp")
	req.Header.Set("X-BRING-CLIENT-SOURCE", "webApp")
	req.Header.Set("X-BRING-COUNTRY", "DE")
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}
	if c.userUUID != "" {
		req.Header.Set("X-BRING-USER-UUID", c.userUUID)
	}
}

// Login authenticates with email and password
func Login(email, password string) (*AuthResponse, error) {
	data := url.Values{}
	data.Set("email", email)
	data.Set("password", password)

	req, err := http.NewRequest("POST", BaseURL+"/v2/bringauth", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-BRING-API-KEY", APIKey)
	req.Header.Set("X-BRING-CLIENT", "webApp")
	req.Header.Set("X-BRING-CLIENT-SOURCE", "webApp")
	req.Header.Set("X-BRING-COUNTRY", "DE")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("login failed (status %d): %s", resp.StatusCode, string(body))
	}

	var auth AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&auth); err != nil {
		return nil, fmt.Errorf("failed to decode auth response: %w", err)
	}
	return &auth, nil
}

// RefreshToken refreshes an expired access token
func RefreshToken(refreshToken string, accessToken, userUUID string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", BaseURL+"/v2/bringauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-BRING-API-KEY", APIKey)
	req.Header.Set("X-BRING-CLIENT", "webApp")
	req.Header.Set("X-BRING-CLIENT-SOURCE", "webApp")
	req.Header.Set("X-BRING-COUNTRY", "DE")
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	if userUUID != "" {
		req.Header.Set("X-BRING-USER-UUID", userUUID)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("refresh failed (status %d): %s", resp.StatusCode, string(body))
	}

	var token TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}
	return &token, nil
}

// GetLists returns all shopping lists for the user
func (c *Client) GetLists() ([]List, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/bringusers/%s/lists", BaseURL, c.userUUID), nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get lists failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get lists failed (status %d): %s", resp.StatusCode, string(body))
	}

	var result ListsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode lists: %w", err)
	}
	return result.Lists, nil
}

// GetItems returns all items on a shopping list
func (c *Client) GetItems(listUUID string) (*ItemsResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/bringlists/%s", BaseURL, listUUID), nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get items failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get items failed (status %d): %s", resp.StatusCode, string(body))
	}

	var result ItemsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode items: %w", err)
	}
	return &result, nil
}

// UpdateItems sends a batch update to a shopping list
func (c *Client) UpdateItems(listUUID string, changes []Change) error {
	body := BatchUpdate{
		Changes: changes,
		Sender:  "",
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v2/bringlists/%s/items", BaseURL, listUUID), bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("update items failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("update items failed (status %d): %s", resp.StatusCode, string(respBody))
	}
	return nil
}

// putForm sends a form-urlencoded PUT to the list endpoint
func (c *Client) putForm(listUUID string, data url.Values) error {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v2/bringlists/%s", BaseURL, listUUID), strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed (status %d): %s", resp.StatusCode, string(body))
	}
	return nil
}

// AddItem adds a single item to a shopping list
func (c *Client) AddItem(listUUID, itemName, spec string) error {
	data := url.Values{}
	data.Set("uuid", listUUID)
	data.Set("purchase", itemName)
	if spec != "" {
		data.Set("specification", spec)
	}
	return c.putForm(listUUID, data)
}

// CompleteItem marks an item as recently bought
func (c *Client) CompleteItem(listUUID, itemName, spec string) error {
	return c.UpdateItems(listUUID, []Change{{
		ItemID:    itemName,
		Spec:      spec,
		Operation: OpToRecently,
	}})
}

// EditItem updates an item's name and/or spec on a shopping list
func (c *Client) EditItem(listUUID, oldName, newName, spec string) error {
	if oldName != newName {
		if err := c.RemoveItem(listUUID, oldName); err != nil {
			return err
		}
	}
	return c.AddItem(listUUID, newName, spec)
}

// RemoveItem removes an item from a shopping list
func (c *Client) RemoveItem(listUUID, itemName string) error {
	return c.UpdateItems(listUUID, []Change{{
		ItemID:    itemName,
		Operation: OpRemove,
	}})
}
