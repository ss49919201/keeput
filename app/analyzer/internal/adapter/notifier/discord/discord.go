package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/samber/mo"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/notifier"
)

type reqBody struct {
	Content string `json:"content"`
}

func loadDiscordWebhookURL() mo.Option[string] {
	url := os.Getenv("DISCORD_WEBHOOK_URL")
	if url == "" {
		return mo.None[string]()
	}
	return mo.Some(url)
}

func createMessage(isGoalAchieved bool) string {
	if isGoalAchieved {
		return "目標達成です🎊よく頑張りました！"
	}
	return "目標未達です😢これから頑張りましょう！"
}

func NewNotify() notifier.Notify {
	return func(req *notifier.NotificationRequest) mo.Result[struct{}] {
		webhookURL := loadDiscordWebhookURL()
		if webhookURL.IsAbsent() {
			return mo.Err[struct{}](fmt.Errorf("discord webhook url is not set or empty"))
		}

		body := reqBody{
			Content: createMessage(req.IsGoalAchieved),
		}

		jsonData, err := json.Marshal(body)
		if err != nil {
			return mo.Err[struct{}](fmt.Errorf("failed to marshal request body: %w", err))
		}

		resp, err := http.Post(webhookURL.MustGet(), "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return mo.Err[struct{}](fmt.Errorf("failed to send request: %w", err))
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return mo.Err[struct{}](fmt.Errorf("failed to request with status code: %d", resp.StatusCode))
		}

		fmt.Println("success request")
		return mo.Ok(struct{}{})
	}
}