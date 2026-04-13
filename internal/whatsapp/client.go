package whatsapp

import (
	"fmt"
	"net/http"
)

const defaultBaseURL = "https://graph.facebook.com/v19.0"

// MessageResponse represents the API response after sending a message.
type MessageResponse struct {
	MessageID string `json:"messages[0].id"`
}

// Client is the WhatsApp Business API client.
type Client struct {
	phoneID     string
	accessToken string
	baseURL     string
	httpClient  *http.Client
}

// NewClient creates a new WhatsApp API client.
// phoneID and accessToken are required.
func NewClient(phoneID, accessToken string, opts ...Option) (*Client, error) {
	if phoneID == "" {
		return nil, fmt.Errorf("phoneID must not be empty")
	}
	if accessToken == "" {
		return nil, fmt.Errorf("accessToken must not be empty")
	}

	c := &Client{
		phoneID:     phoneID,
		accessToken: accessToken,
		baseURL:     defaultBaseURL,
		httpClient:  &http.Client{},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}
