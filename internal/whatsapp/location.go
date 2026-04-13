package whatsapp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// LocationMessage represents a WhatsApp location message payload.
type LocationMessage struct {
	MessagingProduct string   `json:"messaging_product"`
	RecipientType    string   `json:"recipient_type"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Location         Location `json:"location"`
}

// Location holds the geographic coordinates and optional metadata.
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
}

// SendLocationMessage sends a location message to the specified recipient.
func (c *Client) SendLocationMessage(ctx context.Context, to string, lat, lon float64, name, address string) error {
	if to == "" {
		return fmt.Errorf("recipient phone number must not be empty")
	}
	if lat < -90 || lat > 90 {
		return fmt.Errorf("latitude must be between -90 and 90, got %f", lat)
	}
	if lon < -180 || lon > 180 {
		return fmt.Errorf("longitude must be between -180 and 180, got %f", lon)
	}

	payload := LocationMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "location",
		Location: Location{
			Latitude:  lat,
			Longitude: lon,
			Name:      name,
			Address:   address,
		},
	}

	resp, err := c.sendRequest(ctx, http.MethodPost, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiErr map[string]interface{}
		if jsonErr := json.NewDecoder(resp.Body).Decode(&apiErr); jsonErr == nil {
			return fmt.Errorf("API error (status %d): %v", resp.StatusCode, apiErr)
		}
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
