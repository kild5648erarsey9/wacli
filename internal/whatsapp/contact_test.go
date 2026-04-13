package whatsapp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendContactMessage_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewClient("phone-id", "token", WithBaseURL(server.URL))
	contacts := []Contact{
		{Name: ContactName{FormattedName: "Alice Smith", FirstName: "Alice", LastName: "Smith"},
			Phones: []ContactPhone{{Phone: "+1234567890", Type: "MOBILE"}}},
	}
	err := client.SendContactMessage(context.Background(), "+0987654321", contacts)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSendContactMessage_EmptyRecipient(t *testing.T) {
	client, _ := NewClient("phone-id", "token")
	contacts := []Contact{{Name: ContactName{FormattedName: "Alice"}}}
	err := client.SendContactMessage(context.Background(), "", contacts)
	if err == nil {
		t.Fatal("expected error for empty recipient")
	}
}

func TestSendContactMessage_EmptyContacts(t *testing.T) {
	client, _ := NewClient("phone-id", "token")
	err := client.SendContactMessage(context.Background(), "+1234567890", nil)
	if err == nil {
		t.Fatal("expected error for empty contacts slice")
	}
}

func TestSendContactMessage_MissingFormattedName(t *testing.T) {
	client, _ := NewClient("phone-id", "token")
	contacts := []Contact{{Name: ContactName{FormattedName: ""}}}
	err := client.SendContactMessage(context.Background(), "+1234567890", contacts)
	if err == nil {
		t.Fatal("expected error for missing formatted_name")
	}
}

func TestSendContactMessage_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, _ := NewClient("phone-id", "token", WithBaseURL(server.URL))
	contacts := []Contact{{Name: ContactName{FormattedName: "Bob"}}}
	err := client.SendContactMessage(context.Background(), "+1234567890", contacts)
	if err == nil {
		t.Fatal("expected error for server error response")
	}
}
