package internal

import (
	"context"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/ss49919201/fight-op/app/analyzer/internal/model"
)

func Fetch(ctx context.Context, feedURL string) mo.Result[[]*model.Entry] {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		return mo.Err[[]*model.Entry](err)
	}

	return mo.Ok(
		lo.Map(
			feed.Items,
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

func FilterByDateRange(entries []*model.Entry, from, to time.Time) []*model.Entry {
	return lo.Filter(entries, func(entry *model.Entry, _ int) bool {
		return !from.After(entry.PublishedAt) && !to.Before(entry.PublishedAt)
	})
}
