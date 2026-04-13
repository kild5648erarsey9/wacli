package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// TemplateComponent represents a component in a template message.
type TemplateComponent struct {
	Type       string              `json:"type"`
	Parameters []TemplateParameter `json:"parameters"`
}

// TemplateParameter represents a parameter within a template component.
type TemplateParameter struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// templatePayload is the JSON body for sending a template message.
type templatePayload struct {
	MessagingProduct string            `json:"messaging_product"`
	To               string            `json:"to"`
	Type             string            `json:"type"`
	Template         templateBody      `json:"template"`
}

type templateBody struct {
	Name       string               `json:"name"`
	Language   templateLanguage     `json:"language"`
	Components []TemplateComponent  `json:"components,omitempty"`
}

type templateLanguage struct {
	Code string `json:"code"`
}

// SendTemplateMessage sends a WhatsApp template message to the given recipient.
// languageCode defaults to "en_US" if empty.
func (c *Client) SendTemplateMessage(ctx context.Context, recipient, templateName, languageCode string, components []TemplateComponent) error {
	if recipient == "" {
		return fmt.Errorf("recipient phone number must not be empty")
	}
	if templateName == "" {
		return fmt.Errorf("template name must not be empty")
	}
	if languageCode == "" {
		languageCode = "en_US"
	}

	payload := templatePayload{
		MessagingProduct: "whatsapp",
		To:               recipient,
		Type:             "template",
		Template: templateBody{
			Name:       templateName,
			Language:   templateLanguage{Code: languageCode},
			Components: components,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal template payload: %w", err)
	}

	url := fmt.Sprintf("%s/%s/messages", c.baseURL, c.phoneID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send template message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
