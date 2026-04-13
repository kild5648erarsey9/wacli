package whatsapp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendStickerMessage_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"messages":[{"id":"wamid.sticker123"}]}`))
	}))
	defer server.Close()

	client, err := NewClient("test-phone-id", "test-token",
		WithBaseURL(server.URL),
		WithHTTPClient(server.Client()),
	)
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	err = client.SendStickerMessage(context.Background(), "+1234567890", "https://example.com/sticker.webp")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestSendStickerMessage_EmptyRecipient(t *testing.T) {
	client, err := NewClient("test-phone-id", "test-token")
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	err = client.SendStickerMessage(context.Background(), "", "https://example.com/sticker.webp")
	if err == nil {
		t.Fatal("expected error for empty recipient, got nil")
	}
}

func TestSendStickerMessage_EmptyLink(t *testing.T) {
	client, err := NewClient("test-phone-id", "test-token")
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	err = client.SendStickerMessage(context.Background(), "+1234567890", "")
	if err == nil {
		t.Fatal("expected error for empty link, got nil")
	}
}

func TestSendStickerMessage_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"message":"invalid sticker url","code":100}}`))
	}))
	defer server.Close()

	client, err := NewClient("test-phone-id", "test-token",
		WithBaseURL(server.URL),
		WithHTTPClient(server.Client()),
	)
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	err = client.SendStickerMessage(context.Background(), "+1234567890", "https://example.com/sticker.webp")
	if err == nil {
		t.Fatal("expected error for server error response, got nil")
	}
}
