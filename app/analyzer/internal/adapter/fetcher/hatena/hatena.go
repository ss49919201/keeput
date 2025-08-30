package hatena

import (
	"context"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/ss49919201/fight-op/app/analyzer/internal/model"
	"github.com/ss49919201/fight-op/app/analyzer/internal/port/fetcher"
)

func NewFetchAllByDate(feedURL string) fetcher.FetchAllByDate {
	return func(ctx context.Context, from, to time.Time) mo.Result[[]*model.Entry] {
		return FetchAllByDate(ctx, feedURL, from, to)
	}
}

func FetchAllByDate(ctx context.Context, feedURL string, from, to time.Time) mo.Result[[]*model.Entry] {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		return mo.Err[[]*model.Entry](err)
	}

	return mo.Ok(
		lo.Map(
			lo.Filter(feed.Items, func(item *gofeed.Item, _ int) bool {
				return !from.After(*item.PublishedParsed) && !to.Before(*item.PublishedParsed)
			}),
			func(item *gofeed.Item, _ int) *model.Entry {
				return &model.Entry{
					Title:       item.Title,
					Body:        item.Content,
					PublishedAt: *item.PublishedParsed,
				}
			},
		),
	)
}
