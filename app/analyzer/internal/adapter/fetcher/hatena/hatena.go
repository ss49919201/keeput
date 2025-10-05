package hatena

import (
	"context"
	"net/url"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/fetcher/internal"
	"github.com/ss49919201/keeput/app/analyzer/internal/config"
	"github.com/ss49919201/keeput/app/analyzer/internal/model"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/fetcher"
)

func NewFetchLatestEntry() fetcher.FetchLatestEntry {
	return func(ctx context.Context) mo.Result[mo.Option[*model.Entry]] {
		return fetchLatestEntry(ctx, config.FeedURLHatena())
	}
}

func fetchLatestEntry(ctx context.Context, feedURL string) mo.Result[mo.Option[*model.Entry]] {
	u, _ := url.Parse(feedURL)
	q := u.Query()
	q.Set("size", "1")
	u.RawQuery = q.Encode()

	entriesResult := internal.Fetch(ctx, u.String())
	if entriesResult.IsError() {
		return mo.Err[mo.Option[*model.Entry]](entriesResult.Error())
	}

	entries := entriesResult.MustGet()
	latestEntry := lo.MaxBy(entries, func(a *model.Entry, b *model.Entry) bool {
		return a.PublishedAt.After(b.PublishedAt)
	})
	if latestEntry == nil {
		return mo.Ok(mo.None[*model.Entry]())
	}
	return mo.Ok(mo.Some(latestEntry))
}
