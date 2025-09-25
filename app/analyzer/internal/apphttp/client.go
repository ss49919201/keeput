package apphttp

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func DefaultClient() *http.Client {
	client := retryablehttp.NewClient()
	client.RetryMax = 3
	client.Logger = slog.Default()
	standAloneClient := client.StandardClient()
	standAloneClient.Timeout = 30 * time.Second
	standAloneClient.Transport = otelhttp.NewTransport(standAloneClient.Transport)
	return standAloneClient
}
