package fetcher

import (
	"context"
	"time"

	"github.com/samber/mo"
	"github.com/ss49919201/fight-op/app/analyzer/internal/model"
)

type FetchAllByDate = func(ctx context.Context, from, to time.Time) mo.Result[[]*model.Entry]
type FetchLatest = func(context.Context) mo.Result[mo.Option[*model.Entry]]
