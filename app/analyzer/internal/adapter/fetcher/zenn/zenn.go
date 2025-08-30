package zenn

import (
	"context"
	"time"

	"github.com/samber/mo"
	"github.com/ss49919201/fight-op/app/analyzer/internal/adapter/fetcher/internal"
	"github.com/ss49919201/fight-op/app/analyzer/internal/model"
	"github.com/ss49919201/fight-op/app/analyzer/internal/port/fetcher"
)

func NewFetchAllByDate(feedURL string) fetcher.FetchAllByDate {
	return func(ctx context.Context, from, to time.Time) mo.Result[[]*model.Entry] {
		return FetchAllByDate(ctx, feedURL, from, to)
	}
}

func FetchAllByDate(ctx context.Context, feedURL string, from, to time.Time) mo.Result[[]*model.Entry] {
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
