package whatsapp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendReaction_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body reactionMessage
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Type != "reaction" {
			t.Errorf("expected type reaction, got %s", body.Type)
		}
		if body.Reaction.Emoji != "👍" {
			t.Errorf("expected emoji 👍, got %s", body.Reaction.Emoji)
		}
		if body.Reaction.MessageID != "msg-123" {
			t.Errorf("expected messageID msg-123, got %s", body.Reaction.MessageID)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(MessageResponse{MessageID: "resp-1"})
	}))
	defer server.Close()

	client, _ := NewClient("phone-id", "token", WithBaseURL(server.URL))
	resp, err := client.SendReaction(context.Background(), "+1234567890", "msg-123", "👍")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MessageID != "resp-1" {
		t.Errorf("expected resp-1, got %s", resp.MessageID)
	}
}

func TestSendReaction_RemoveEmoji(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body reactionMessage
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body.Reaction.Emoji != "" {
			t.Errorf("expected empty emoji for removal, got %s", body.Reaction.Emoji)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(MessageResponse{MessageID: "resp-2"})
	}))
	defer server.Close()

	client, _ := NewClient("phone-id", "token", WithBaseURL(server.URL))
	_, err := client.SendReaction(context.Background(), "+1234567890", "msg-123", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSendReaction_EmptyRecipient(t *testing.T) {
	client, _ := NewClient("phone-id", "token")
	_, err := client.SendReaction(context.Background(), "", "msg-123", "👍")
	if err == nil {
		t.Fatal("expected error for empty recipient")
	}
}

func TestSendReaction_EmptyMessageID(t *testing.T) {
	client, _ := NewClient("phone-id", "token")
	_, err := client.SendReaction(context.Background(), "+1234567890", "", "👍")
	if err == nil {
		t.Fatal("expected error for empty messageID")
	}
}

func TestSendReaction_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client, _ := NewClient("phone-id", "token", WithBaseURL(server.URL))
	_, err := client.SendReaction(context.Background(), "+1234567890", "msg-123", "👍")
	if err == nil {
		t.Fatal("expected error for server error response")
	}
}
