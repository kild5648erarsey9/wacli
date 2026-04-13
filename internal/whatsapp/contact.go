package whatsapp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ContactName represents the name fields of a contact.
type ContactName struct {
	FormattedName string `json:"formatted_name"`
	FirstName     string `json:"first_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
}

// ContactPhone represents a phone entry for a contact.
type ContactPhone struct {
	Phone string `json:"phone"`
	Type  string `json:"type,omitempty"`
}

// Contact represents a single contact card.
type Contact struct {
	Name   ContactName    `json:"name"`
	Phones []ContactPhone `json:"phones,omitempty"`
}

type contactMessage struct {
	MessagingProduct string    `json:"messaging_product"`
	To               string    `json:"to"`
	Type             string    `json:"type"`
	Contacts         []Contact `json:"contacts"`
}

// SendContactMessage sends one or more contact cards to the given recipient.
func (c *Client) SendContactMessage(ctx context.Context, to string, contacts []Contact) error {
	if strings.TrimSpace(to) == "" {
		return fmt.Errorf("recipient phone number must not be empty")
	}
	if len(contacts) == 0 {
		return fmt.Errorf("at least one contact must be provided")
	}
	for i, ct := range contacts {
		if strings.TrimSpace(ct.Name.FormattedName) == "" {
			return fmt.Errorf("contact[%d]: formatted_name must not be empty", i)
		}
	}

	payload := contactMessage{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "contacts",
		Contacts:         contacts,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal contact message: %w", err)
	}

	resp, err := c.post(ctx, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
