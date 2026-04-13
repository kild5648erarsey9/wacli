package whatsapp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendLocationMessage_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"messages":[{"id":"msg-id-123"}]}`))
	}))
	defer server.Close()

	client, err := NewClient("phone-id", "token", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	err = client.SendLocationMessage(context.Background(), "+1234567890", 37.7749, -122.4194, "San Francisco", "CA, USA")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestSendLocationMessage_EmptyRecipient(t *testing.T) {
	client, err := NewClient("phone-id", "token")
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	err = client.SendLocationMessage(context.Background(), "", 37.7749, -122.4194, "", "")
	if err == nil {
		t.Fatal("expected error for empty recipient, got nil")
	}
}

func TestSendLocationMessage_InvalidLatitude(t *testing.T) {
	client, err := NewClient("phone-id", "token")
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	err = client.SendLocationMessage(context.Background(), "+1234567890", 95.0, -122.4194, "", "")
	if err == nil {
		t.Fatal("expected error for invalid latitude, got nil")
	}
}

func TestSendLocationMessage_InvalidLongitude(t *testing.T) {
	client, err := NewClient("phone-id", "token")
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	err = client.SendLocationMessage(context.Background(), "+1234567890", 37.7749, 200.0, "", "")
	if err == nil {
		t.Fatal("expected error for invalid longitude, got nil")
	}
}

func TestSendLocationMessage_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":{"message":"internal server error"}}`))
	}))
	defer server.Close()

	client, err := NewClient("phone-id", "token", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	err = client.SendLocationMessage(context.Background(), "+1234567890", 37.7749, -122.4194, "", "")
	if err == nil {
		t.Fatal("expected error from server, got nil")
	}
}

func TestSendLocationMessage_NoNameOrAddress(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"messages":[{"id":"msg-id-456"}]}`))
	}))
	defer server.Close()

	client, err := NewClient("phone-id", "token", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	err = client.SendLocationMessage(context.Background(), "+1234567890", 0.0, 0.0, "", "")
	if err != nil {
		t.Fatalf("expected no error for null island coordinates, got: %v", err)
	}
}
