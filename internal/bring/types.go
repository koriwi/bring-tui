package bring

const (
	BaseURL = "https://api.getbring.com/rest"
	APIKey  = "cof4Nc6D8saplXjE3h3HXqHH8m7VU2i1Gs0g85Sp"

	OpToPurchase = "TO_PURCHASE"
	OpToRecently = "TO_RECENTLY"
	OpRemove     = "REMOVE"
)

// AuthResponse is returned from POST /v2/bringauth
type AuthResponse struct {
	UUID          string `json:"uuid"`
	PublicUUID    string `json:"publicUuid"`
	BringListUUID string `json:"bringListUUID"`
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	ExpiresIn     int    `json:"expires_in"`
	TokenType     string `json:"token_type"`
	Email         string `json:"email"`
	Name          string `json:"name"`
}

// TokenResponse is returned from POST /v2/bringauth/token (refresh)
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// ListsResponse is returned from GET /bringusers/{uuid}/lists
type ListsResponse struct {
	Lists []List `json:"lists"`
}

// List represents a single shopping list
type List struct {
	ListUUID string `json:"listUuid"`
	Name     string `json:"name"`
	Theme    string `json:"theme"`
}

// ItemsResponse is returned from GET /v2/bringlists/{listUuid}
type ItemsResponse struct {
	UUID   string    `json:"uuid"`
	Status string    `json:"status"`
	Items  ItemGroup `json:"items"`
}

// ItemGroup contains active and recently bought items
type ItemGroup struct {
	Purchase []Item `json:"purchase"`
	Recently []Item `json:"recently"`
}

// Item represents a single shopping list item
type Item struct {
	UUID   string `json:"uuid"`
	ItemID string `json:"itemId"`
	Spec   string `json:"specification"`
}

// BatchUpdate is the request body for PUT /v2/bringlists/{listUuid}/items
type BatchUpdate struct {
	Changes []Change `json:"changes"`
	Sender  string   `json:"sender"`
}

// Change represents a single item operation
type Change struct {
	ItemID    string `json:"itemId"`
	Spec      string `json:"spec"`
	Operation string `json:"operation"`
}

