package slack

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "log/slog"
    "net/http"
    "sync"

    "github.com/hashicorp/go-retryablehttp"
    "github.com/ss49919201/keeput/app/analyzer/internal/config"
    "github.com/ss49919201/keeput/app/analyzer/internal/model"
    "github.com/ss49919201/keeput/app/analyzer/internal/port/notifier"
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var httpClient = sync.OnceValue(func() *http.Client {
    client := retryablehttp.NewClient()
    client.RetryMax = 3
    client.Logger = slog.Default()

    std := client.StandardClient()
    std.Transport = otelhttp.NewTransport(std.Transport)
    return std
})

type slackReqBody struct {
    Text string `json:"text"`
}

// NewNotify builds a notifier.Notify that posts to Slack Incoming Webhook.
// It uses SLACK_WEBHOOK_URL from environment via config.
func NewNotify() notifier.Notify {
    return func(ctx context.Context, report *model.AnalysisReport) error {
        url := config.SlackWebhookURL()
        if url == "" {
            return fmt.Errorf("SLACK_WEBHOOK_URL is empty")
        }
        msg := "目標未達です😢これから頑張りましょう！"
        if report.IsGoalAchieved {
            msg = "目標達成です🎊よく頑張りました！"
        }
        body := slackReqBody{Text: msg}
        b, err := json.Marshal(body)
        if err != nil {
            return err
        }
        req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(b))
        if err != nil {
            return err
        }
        req.Header.Set("Content-Type", "application/json")
        resp, err := httpClient().Do(req)
        if err != nil {
            return err
        }
        defer resp.Body.Close()
        if resp.StatusCode != http.StatusOK {
            rb, _ := io.ReadAll(resp.Body)
            return fmt.Errorf("slack webhook unexpected status code %d, body %s", resp.StatusCode, string(rb))
        }
        return nil
    }
}
