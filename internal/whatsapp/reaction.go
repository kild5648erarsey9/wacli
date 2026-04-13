package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type reactionMessage struct {
	MessagingProduct string   `json:"messaging_product"`
	RecipientType    string   `json:"recipient_type"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Reaction         reaction `json:"reaction"`
}

type reaction struct {
	MessageID string `json:"message_id"`
	Emoji     string `json:"emoji"`
}

// SendReaction sends a reaction (emoji) to a specific message.
// Pass an empty emoji string to remove an existing reaction.
func (c *Client) SendReaction(ctx context.Context, recipient, messageID, emoji string) (*MessageResponse, error) {
	if recipient == "" {
		return nil, fmt.Errorf("recipient must not be empty")
	}
	if messageID == "" {
		return nil, fmt.Errorf("messageID must not be empty")
	}

	body := reactionMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               recipient,
		Type:             "reaction",
		Reaction: reaction{
			MessageID: messageID,
			Emoji:     emoji,
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal reaction message: %w", err)
	}

	url := fmt.Sprintf("%s/%s/messages", c.baseURL, c.phoneID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send reaction: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result MessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}
