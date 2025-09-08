package internal

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/mmcdole/gofeed"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/ss49919201/keeput/app/analyzer/internal/model"
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

// NOTE: 公開日が存在しないエントリは除外する。
func Fetch(ctx context.Context, feedURL string) mo.Result[[]*model.Entry] {
	fp := gofeed.NewParser()
	fp.Client = httpClient()
	feed, err := fp.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		return mo.Err[[]*model.Entry](err)
	}

	return mo.Ok(
		lo.FilterMap(
			feed.Items,
			func(item *gofeed.Item, _ int) (*model.Entry, bool) {
				if item.PublishedParsed == nil {
					return nil, false
				}

				return &model.Entry{
					Title:       item.Title,
					Body:        item.Content,
					PublishedAt: *item.PublishedParsed,
				}, true
			},
		),
	)
}

func FilterByDateRange(entries []*model.Entry, from, to time.Time) []*model.Entry {
	return lo.Filter(entries, func(entry *model.Entry, _ int) bool {
		return !from.After(entry.PublishedAt) && !to.Before(entry.PublishedAt)
	})
}
