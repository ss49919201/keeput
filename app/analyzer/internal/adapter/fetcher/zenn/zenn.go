package zenn

import (
	"context"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/fetcher/internal"
	"github.com/ss49919201/keeput/app/analyzer/internal/config"
	"github.com/ss49919201/keeput/app/analyzer/internal/model"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/fetcher"
)

func NewFetchLatest() fetcher.FetchLatest {
	return func(ctx context.Context) mo.Result[mo.Option[*model.Entry]] {
		return fetchLatest(ctx, config.FeedURLZenn())
	}
}

func fetchLatest(ctx context.Context, feedURL string) mo.Result[mo.Option[*model.Entry]] {
	entriesResult := internal.Fetch(ctx, feedURL)
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
