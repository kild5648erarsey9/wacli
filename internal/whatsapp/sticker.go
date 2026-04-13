package whatsapp

import (
	"context"
	"fmt"
	"net/http"
)

// StickerMessage represents a WhatsApp sticker message payload.
type StickerMessage struct {
	MessagingProduct string          `json:"messaging_product"`
	RecipientType    string          `json:"recipient_type"`
	To               string          `json:"to"`
	Type             string          `json:"type"`
	Sticker          StickerObject   `json:"sticker"`
}

// StickerObject holds the sticker link.
type StickerObject struct {
	Link string `json:"link"`
}

// SendStickerMessage sends a sticker message to the specified recipient.
// The link must point to a WebP image (static or animated).
func (c *Client) SendStickerMessage(ctx context.Context, recipient, link string) error {
	if recipient == "" {
		return fmt.Errorf("recipient phone number must not be empty")
	}
	if link == "" {
		return fmt.Errorf("sticker link must not be empty")
	}

	payload := StickerMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               recipient,
		Type:             "sticker",
		Sticker: StickerObject{
			Link: link,
		},
	}

	resp, err := c.sendRequest(ctx, http.MethodPost, payload)
	if err != nil {
		return fmt.Errorf("sending sticker message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return parseErrorResponse(resp)
	}
	return nil
}
