package whatsapp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// InteractiveButton represents a reply button in an interactive message.
type InteractiveButton struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// InteractiveMessageRequest is the payload for sending an interactive button message.
type InteractiveMessageRequest struct {
	MessagingProduct string              `json:"messaging_product"`
	RecipientType    string              `json:"recipient_type"`
	To               string              `json:"to"`
	Type             string              `json:"type"`
	Interactive      InteractivePayload  `json:"interactive"`
}

// InteractivePayload holds the interactive message body and action.
type InteractivePayload struct {
	Type   string            `json:"type"`
	Body   InteractiveBody   `json:"body"`
	Action InteractiveAction `json:"action"`
}

// InteractiveBody is the text body of the interactive message.
type InteractiveBody struct {
	Text string `json:"text"`
}

// InteractiveAction contains the list of reply buttons.
type InteractiveAction struct {
	Buttons []InteractiveButtonWrapper `json:"buttons"`
}

// InteractiveButtonWrapper wraps a reply button for the API payload.
type InteractiveButtonWrapper struct {
	Type  string             `json:"type"`
	Reply InteractiveButton  `json:"reply"`
}

// SendInteractiveButtonMessage sends an interactive reply-button message via the WhatsApp API.
func (c *Client) SendInteractiveButtonMessage(ctx context.Context, to, bodyText string, buttons []InteractiveButton) error {
	if to == "" {
		return fmt.Errorf("recipient phone number must not be empty")
	}
	if bodyText == "" {
		return fmt.Errorf("body text must not be empty")
	}
	if len(buttons) == 0 {
		return fmt.Errorf("at least one button is required")
	}
	if len(buttons) > 3 {
		return fmt.Errorf("interactive messages support at most 3 buttons, got %d", len(buttons))
	}

	wrappers := make([]InteractiveButtonWrapper, len(buttons))
	for i, b := range buttons {
		if b.ID == "" || b.Title == "" {
			return fmt.Errorf("button at index %d must have a non-empty id and title", i)
		}
		wrappers[i] = InteractiveButtonWrapper{Type: "reply", Reply: b}
	}

	payload := InteractiveMessageRequest{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "interactive",
		Interactive: InteractivePayload{
			Type: "button",
			Body: InteractiveBody{Text: bodyText},
			Action: InteractiveAction{Buttons: wrappers},
		},
	}

	resp, err := c.sendRequest(ctx, http.MethodPost, c.messagesURL(), payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if jsonErr := json.NewDecoder(resp.Body).Decode(&apiErr); jsonErr == nil && apiErr.Error.Message != "" {
			return fmt.Errorf("whatsapp api error: %s", apiErr.Error.Message)
		}
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
