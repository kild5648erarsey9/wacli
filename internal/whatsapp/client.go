package whatsapp

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "https://graph.facebook.com/v18.0"
	defaultTimeout = 30 * time.Second
)

// Client holds configuration and HTTP client for WhatsApp API calls.
type Client struct {
	baseURL    string
	phoneID    string
	accessToken string
	httpClient *http.Client
}

// Config holds the configuration for creating a new Client.
type Config struct {
	PhoneID     string
	AccessToken string
	BaseURL     string
	Timeout     time.Duration
}

// NewClient creates a new WhatsApp API client from the given config.
func NewClient(cfg Config) (*Client, error) {
	if cfg.PhoneID == "" {
		return nil, fmt.Errorf("phone ID must not be empty")
	}
	if cfg.AccessToken == "" {
		return nil, fmt.Errorf("access token must not be empty")
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}

	return &Client{
		baseURL:     baseURL,
		phoneID:     cfg.PhoneID,
		accessToken: cfg.AccessToken,
		httpClient:  &http.Client{Timeout: timeout},
	}, nil
}

// SendTextMessage sends a plain text message to the given recipient.
func (c *Client) SendTextMessage(ctx context.Context, to, body string) error {
	if to == "" {
		return fmt.Errorf("recipient phone number must not be empty")
	}
	if body == "" {
		return fmt.Errorf("message body must not be empty")
	}

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "text",
		"text":              map[string]string{"body": body},
	}

	url := fmt.Sprintf("%s/%s/messages", c.baseURL, c.phoneID)
	return c.doPost(ctx, url, payload)
}
