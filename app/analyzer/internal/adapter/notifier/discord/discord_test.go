package discord

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

type payload struct {
    Content string `json:"content"`
}

func TestNotify_Success_OnAchieved(t *testing.T) {
    var received payload
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer r.Body.Close()
        if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
            t.Fatalf("failed to decode: %v", err)
        }
        w.WriteHeader(http.StatusOK)
    }))
    t.Cleanup(ts.Close)

    err := notify(context.Background(), ts.URL, true)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if received.Content == "" {
        t.Fatalf("expected content, got empty")
    }
}

func TestNotify_Success_OnNotAchieved(t *testing.T) {
    var received payload
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer r.Body.Close()
        if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
            t.Fatalf("failed to decode: %v", err)
        }
        w.WriteHeader(http.StatusOK)
    }))
    t.Cleanup(ts.Close)

    err := notify(context.Background(), ts.URL, false)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if received.Content == "" {
        t.Fatalf("expected content, got empty")
    }
}

func TestNotify_ErrorWhenWebhookMissing(t *testing.T) {
    if err := notify(context.Background(), "", true); err == nil {
        t.Fatalf("expected error for missing webhook url")
    }
}

func TestNotify_ErrorOnNon200(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusBadRequest)
        _, _ = w.Write([]byte(`{"error":"bad"}`))
    }))
    t.Cleanup(ts.Close)

    if err := notify(context.Background(), ts.URL, true); err == nil {
        t.Fatalf("expected error on non-200 response")
    }
}

