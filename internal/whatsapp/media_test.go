package whatsapp

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendMediaMessage_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing or wrong Authorization header")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewClient("phone-id", "test-token", WithBaseURL(server.URL))
	err := client.SendMediaMessage("+1234567890", MediaTypeImage, "https://example.com/image.jpg", "Hello")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSendMediaMessage_EmptyRecipient(t *testing.T) {
	client, _ := NewClient("phone-id", "test-token")
	err := client.SendMediaMessage("", MediaTypeImage, "https://example.com/image.jpg", "")
	if err == nil {
		t.Fatal("expected error for empty recipient")
	}
}

func TestSendMediaMessage_EmptyLink(t *testing.T) {
	client, _ := NewClient("phone-id", "test-token")
	err := client.SendMediaMessage("+1234567890", MediaTypeImage, "", "")
	if err == nil {
		t.Fatal("expected error for empty link")
	}
}

func TestSendMediaMessage_UnsupportedType(t *testing.T) {
	client, _ := NewClient("phone-id", "test-token")
	err := client.SendMediaMessage("+1234567890", MediaType("sticker"), "https://example.com/s.webp", "")
	if err == nil {
		t.Fatal("expected error for unsupported media type")
	}
}

func TestSendMediaMessage_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, _ := NewClient("phone-id", "test-token", WithBaseURL(server.URL))
	err := client.SendMediaMessage("+1234567890", MediaTypeDocument, "https://example.com/doc.pdf", "My Doc")
	if err == nil {
		t.Fatal("expected error on server error response")
	}
}

func TestSendMediaMessage_AllMediaTypes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewClient("phone-id", "test-token", WithBaseURL(server.URL))

	types := []MediaType{MediaTypeImage, MediaTypeDocument, MediaTypeAudio, MediaTypeVideo}
	for _, mt := range types {
		err := client.SendMediaMessage("+1234567890", mt, "https://example.com/file", "caption")
		if err != nil {
			t.Errorf("expected no error for media type %s, got %v", mt, err)
		}
	}
}
