package whatsapp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient_MissingPhoneID(t *testing.T) {
	_, err := NewClient(Config{AccessToken: "token"})
	if err == nil {
		t.Fatal("expected error for missing phone ID")
	}
}

func TestNewClient_MissingAccessToken(t *testing.T) {
	_, err := NewClient(Config{PhoneID: "123"})
	if err == nil {
		t.Fatal("expected error for missing access token")
	}
}

func TestNewClient_Defaults(t *testing.T) {
	c, err := NewClient(Config{PhoneID: "123", AccessToken: "token"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.baseURL != defaultBaseURL {
		t.Errorf("expected default base URL %q, got %q", defaultBaseURL, c.baseURL)
	}
	if c.httpClient.Timeout != defaultTimeout {
		t.Errorf("expected default timeout %v, got %v", defaultTimeout, c.httpClient.Timeout)
	}
}

func TestSendTextMessage_Success(t *testing.T) {
	var received map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"messages":[{"id":"wamid.test"}]}`))
	}))
	defer ts.Close()

	c, _ := NewClient(Config{
		PhoneID:     "phone123",
		AccessToken: "token",
		BaseURL:     ts.URL,
	})

	if err := c.SendTextMessage(context.Background(), "+1234567890", "Hello!"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["to"] != "+1234567890" {
		t.Errorf("expected to=+1234567890, got %v", received["to"])
	}
}

func TestSendTextMessage_EmptyRecipient(t *testing.T) {
	c, _ := NewClient(Config{PhoneID: "123", AccessToken: "token"})
	err := c.SendTextMessage(context.Background(), "", "Hello")
	if err == nil {
		t.Fatal("expected error for empty recipient")
	}
}

func TestSendTextMessage_EmptyBody(t *testing.T) {
	c, _ := NewClient(Config{PhoneID: "123", AccessToken: "token"})
	err := c.SendTextMessage(context.Background(), "+1234567890", "")
	if err == nil {
		t.Fatal("expected error for empty body")
	}
}
