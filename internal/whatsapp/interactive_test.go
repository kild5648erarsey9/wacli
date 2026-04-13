package whatsapp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendInteractiveButtonMessage_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload InteractiveMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if payload.To != "+1234567890" {
			t.Errorf("expected To=+1234567890, got %s", payload.To)
		}
		if payload.Interactive.Body.Text != "Pick an option" {
			t.Errorf("unexpected body text: %s", payload.Interactive.Body.Text)
		}
		if len(payload.Interactive.Action.Buttons) != 2 {
			t.Errorf("expected 2 buttons, got %d", len(payload.Interactive.Action.Buttons))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewClient("phone-id", "token", WithBaseURL(server.URL))
	buttons := []InteractiveButton{
		{ID: "btn1", Title: "Yes"},
		{ID: "btn2", Title: "No"},
	}
	err := client.SendInteractiveButtonMessage(context.Background(), "+1234567890", "Pick an option", buttons)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSendInteractiveButtonMessage_EmptyRecipient(t *testing.T) {
	client, _ := NewClient("phone-id", "token")
	err := client.SendInteractiveButtonMessage(context.Background(), "", "body", []InteractiveButton{{ID: "1", Title: "A"}})
	if err == nil {
		t.Fatal("expected error for empty recipient")
	}
}

func TestSendInteractiveButtonMessage_EmptyBodyText(t *testing.T) {
	client, _ := NewClient("phone-id", "token")
	err := client.SendInteractiveButtonMessage(context.Background(), "+1234567890", "", []InteractiveButton{{ID: "1", Title: "A"}})
	if err == nil {
		t.Fatal("expected error for empty body text")
	}
}

func TestSendInteractiveButtonMessage_NoButtons(t *testing.T) {
	client, _ := NewClient("phone-id", "token")
	err := client.SendInteractiveButtonMessage(context.Background(), "+1234567890", "body", nil)
	if err == nil {
		t.Fatal("expected error for no buttons")
	}
}

func TestSendInteractiveButtonMessage_TooManyButtons(t *testing.T) {
	client, _ := NewClient("phone-id", "token")
	buttons := []InteractiveButton{
		{ID: "1", Title: "A"},
		{ID: "2", Title: "B"},
		{ID: "3", Title: "C"},
		{ID: "4", Title: "D"},
	}
	err := client.SendInteractiveButtonMessage(context.Background(), "+1234567890", "body", buttons)
	if err == nil {
		t.Fatal("expected error for more than 3 buttons")
	}
}

func TestSendInteractiveButtonMessage_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(APIError{Error: APIErrorDetail{Message: "internal server error"}})
	}))
	defer server.Close()

	client, _ := NewClient("phone-id", "token", WithBaseURL(server.URL))
	buttons := []InteractiveButton{{ID: "1", Title: "A"}}
	err := client.SendInteractiveButtonMessage(context.Background(), "+1234567890", "body", buttons)
	if err == nil {
		t.Fatal("expected error on server error response")
	}
}
