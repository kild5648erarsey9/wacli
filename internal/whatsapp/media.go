package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// MediaType represents the type of media to send.
type MediaType string

const (
	MediaTypeImage    MediaType = "image"
	MediaTypeDocument MediaType = "document"
	MediaTypeAudio    MediaType = "audio"
	MediaTypeVideo    MediaType = "video"
)

// mediaMessage represents the payload for sending a media message.
type mediaMessage struct {
	MessagingProduct string    `json:"messaging_product"`
	RecipientType    string    `json:"recipient_type"`
	To               string    `json:"to"`
	Type             MediaType `json:"type"`
	Image            *mediaObject `json:"image,omitempty"`
	Document         *mediaObject `json:"document,omitempty"`
	Audio            *mediaObject `json:"audio,omitempty"`
	Video            *mediaObject `json:"video,omitempty"`
}

// mediaObject holds the link and optional caption for a media item.
type mediaObject struct {
	Link    string `json:"link"`
	Caption string `json:"caption,omitempty"`
}

// SendMediaMessage sends a media message (image, document, audio, video) to the given recipient.
func (c *Client) SendMediaMessage(recipient string, mediaType MediaType, link string, caption string) error {
	if recipient == "" {
		return fmt.Errorf("recipient phone number must not be empty")
	}
	if link == "" {
		return fmt.Errorf("media link must not be empty")
	}

	obj := &mediaObject{Link: link, Caption: caption}

	msg := mediaMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               recipient,
		Type:             mediaType,
	}

	switch mediaType {
	case MediaTypeImage:
		msg.Image = obj
	case MediaTypeDocument:
		msg.Document = obj
	case MediaTypeAudio:
		msg.Audio = obj
	case MediaTypeVideo:
		msg.Video = obj
	default:
		return fmt.Errorf("unsupported media type: %s", mediaType)
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal media message: %w", err)
	}

	url := fmt.Sprintf("%s/%s/messages", c.baseURL, c.phoneID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send media message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
