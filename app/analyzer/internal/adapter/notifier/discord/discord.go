package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ss49919201/keeput/app/analyzer/internal/config"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/notifier"
)

type reqBody struct {
	Content string `json:"content"`
}

func NewNotifyAnalysisReport() notifier.NotifyAnalysisReport {
	return func(ctx context.Context, isGoalAchieved bool) error {
		return notify(ctx, config.DiscordWebhookURL(), isGoalAchieved)
	}
}

func notify(ctx context.Context, webhookURL string, isGoalAchieved bool) error {
	body := reqBody{Content: message(isGoalAchieved)}
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to request with status code: %d; body: %s", resp.StatusCode, string(b))
	}

	return nil
}

func message(isGoalAchieved bool) string {
	if isGoalAchieved {
		return "目標達成です🎊よく頑張りました！"
	}
	return "目標未達です😢これから頑張りましょう！"
}
