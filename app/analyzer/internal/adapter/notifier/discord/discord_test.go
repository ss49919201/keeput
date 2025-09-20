package discord

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ss49919201/keeput/app/analyzer/internal/port/notifier"
)

func TestCreateMessage(t *testing.T) {
	tests := []struct {
		name           string
		isGoalAchieved bool
		expected       string
	}{
		{
			name:           "goal achieved",
			isGoalAchieved: true,
			expected:       "目標達成です🎊よく頑張りました！",
		},
		{
			name:           "goal not achieved",
			isGoalAchieved: false,
			expected:       "目標未達です😢これから頑張りましょう！",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := createMessage(tt.isGoalAchieved)
			if result != tt.expected {
				t.Errorf("createMessage() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLoadDiscordWebhookURL(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		hasValue bool
	}{
		{
			name:     "url is set",
			envValue: "https://discord.com/api/webhooks/123/456",
			hasValue: true,
		},
		{
			name:     "url is empty",
			envValue: "",
			hasValue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			os.Setenv("DISCORD_WEBHOOK_URL", tt.envValue)
			defer os.Unsetenv("DISCORD_WEBHOOK_URL")

			result := loadDiscordWebhookURL()
			if result.IsPresent() != tt.hasValue {
				t.Errorf("loadDiscordWebhookURL() presence = %v, want %v", result.IsPresent(), tt.hasValue)
			}

			if tt.hasValue && result.MustGet() != tt.envValue {
				t.Errorf("loadDiscordWebhookURL() = %v, want %v", result.MustGet(), tt.envValue)
			}
		})
	}
}

func TestNewNotify(t *testing.T) {
	tests := []struct {
		name           string
		webhookURL     string
		isGoalAchieved bool
		statusCode     int
		expectError    bool
	}{
		{
			name:           "successful notification - goal achieved",
			webhookURL:     "valid",
			isGoalAchieved: true,
			statusCode:     200,
			expectError:    false,
		},
		{
			name:           "successful notification - goal not achieved",
			webhookURL:     "valid",
			isGoalAchieved: false,
			statusCode:     200,
			expectError:    false,
		},
		{
			name:           "failed request - bad status code",
			webhookURL:     "valid",
			isGoalAchieved: true,
			statusCode:     400,
			expectError:    true,
		},
		{
			name:           "missing webhook url",
			webhookURL:     "",
			isGoalAchieved: true,
			statusCode:     200,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock server
			var server *httptest.Server
			if tt.webhookURL == "valid" {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.statusCode)
				}))
				defer server.Close()
				os.Setenv("DISCORD_WEBHOOK_URL", server.URL)
			} else {
				os.Setenv("DISCORD_WEBHOOK_URL", tt.webhookURL)
			}
			defer os.Unsetenv("DISCORD_WEBHOOK_URL")

			notify := NewNotify()
			result := notify(&notifier.NotificationRequest{
				IsGoalAchieved: tt.isGoalAchieved,
			})

			if tt.expectError && result.IsOk() {
				t.Error("expected error but got success")
			}
			if !tt.expectError && result.IsError() {
				t.Errorf("expected success but got error: %v", result.Error())
			}
		})
	}
}