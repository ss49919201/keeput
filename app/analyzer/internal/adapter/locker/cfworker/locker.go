package cfworker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"sync"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/samber/mo"
	"github.com/ss49919201/keeput/app/analyzer/internal/config"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/locker"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var httpClient = sync.OnceValue(func() *http.Client {
	// NOTE: retryablehttp.NewClient() は内部で cleanhttp.DefaultPooledClient() を使う。
	// cleanhttp.DefaultPooledClient() が返す http.Client にはタイムアウトが設定されている。
	client := retryablehttp.NewClient()
	client.RetryMax = 3
	client.Logger = slog.Default()

	standAloneClient := client.StandardClient()
	standAloneClient.Transport = otelhttp.NewTransport(standAloneClient.Transport)
	return standAloneClient
})

func NewAcquire() locker.Acquire {
	return func(ctx context.Context, lockID string) mo.Result[bool] {
		return acquire(ctx, lockID, config.LockerURLCloudflareWorker())
	}
}

type acquireResponse struct {
	Msg string `json:"msg"`
}

// TODO: release との共通部分を抽象化して関数にする。
func acquire(ctx context.Context, lockID, baseURL string) mo.Result[bool] {
	reqBodyMap := map[string]string{
		"lockId": lockID,
	}
	reqBodyBytes, err := json.Marshal(reqBodyMap)
	if err != nil {
		return mo.Err[bool](err)
	}
	url, err := url.JoinPath(baseURL, "/acquire")
	if err != nil {
		return mo.Err[bool](err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return mo.Err[bool](err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient().Do(req)
	if err != nil {
		return mo.Err[bool](err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return mo.Err[bool](fmt.Errorf("failed to acquire lock, unexpected status code %d, body %s", resp.StatusCode, body))
	}
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return mo.Err[bool](err)
	}
	var unmarshaledResp *acquireResponse
	if err := json.Unmarshal(respBodyBytes, &unmarshaledResp); err != nil {
		return mo.Err[bool](err)
	}
	return mo.Ok(unmarshaledResp.Msg == "ok")
}

func NewRelease() locker.Release {
	return func(ctx context.Context, lockID string) error {
		return release(ctx, lockID, config.LockerURLCloudflareWorker())
	}
}

type releaseResponse struct {
	Msg string `json:"msg"`
}

func release(ctx context.Context, lockID, baseURL string) error {
	reqBodyMap := map[string]string{
		"lockId": lockID,
	}
	reqBodyBytes, err := json.Marshal(reqBodyMap)
	if err != nil {
		return err
	}
	url, err := url.JoinPath(baseURL, "/release")
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to release lock, unexpected status code %d, body %s", resp.StatusCode, body)
	}
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var unmarshaledResp *releaseResponse
	if err := json.Unmarshal(respBodyBytes, &unmarshaledResp); err != nil {
		return err
	}
	if unmarshaledResp.Msg != "ok" {
		return errors.New("failed release lock")
	}
	return nil
}
