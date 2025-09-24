package hatena

import (
	"context"
	"net/url"
	"time"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/fetcher/internal"
	"github.com/ss49919201/keeput/app/analyzer/internal/config"
	"github.com/ss49919201/keeput/app/analyzer/internal/model"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/fetcher"
)

func NewFetchAllByDate() fetcher.FetchAllByDate {
	return func(ctx context.Context, from, to time.Time) mo.Result[[]*model.Entry] {
		return fetchAllByDate(ctx, config.FeedURLHatena(), from, to)
	}
}

func fetchAllByDate(ctx context.Context, feedURL string, from, to time.Time) mo.Result[[]*model.Entry] {
	entries := internal.Fetch(ctx, feedURL)
	return entries.Match(
		func(entries []*model.Entry) ([]*model.Entry, error) {
			return internal.FilterByDateRange(entries, from, to), nil
		},
		func(err error) ([]*model.Entry, error) {
			return nil, err
		},
	)
}

func NewFetchLatest() fetcher.FetchLatest {
	return func(ctx context.Context) mo.Result[mo.Option[*model.Entry]] {
		return fetchLatest(ctx, config.FeedURLHatena())
	}
}

func fetchLatest(ctx context.Context, feedURL string) mo.Result[mo.Option[*model.Entry]] {
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
