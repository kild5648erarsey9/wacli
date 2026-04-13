package whatsapp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendTemplateMessage_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("unexpected Authorization header: %s", r.Header.Get("Authorization"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, err := NewClient("phone-123", "test-token", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	err = c.SendTemplateMessage(context.Background(), "+1234567890", "hello_world", "en_US", nil)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestSendTemplateMessage_DefaultLanguage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, err := NewClient("phone-123", "test-token", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	// Empty language code should default to en_US without error
	err = c.SendTemplateMessage(context.Background(), "+1234567890", "hello_world", "", nil)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestSendTemplateMessage_EmptyRecipient(t *testing.T) {
	c, err := NewClient("phone-123", "test-token")
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	err = c.SendTemplateMessage(context.Background(), "", "hello_world", "en_US", nil)
	if err == nil {
		t.Fatal("expected error for empty recipient, got nil")
	}
}

func TestSendTemplateMessage_EmptyTemplateName(t *testing.T) {
	c, err := NewClient("phone-123", "test-token")
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	err = c.SendTemplateMessage(context.Background(), "+1234567890", "", "en_US", nil)
	if err == nil {
		t.Fatal("expected error for empty template name, got nil")
	}
}

func TestSendTemplateMessage_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c, err := NewClient("phone-123", "test-token", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	err = c.SendTemplateMessage(context.Background(), "+1234567890", "hello_world", "en_US", nil)
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}
