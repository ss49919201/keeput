package fetcher

import (
	"context"

	"github.com/samber/mo"
	"github.com/ss49919201/keeput/app/analyzer/internal/model"
)

type FetchLatestEntry = func(context.Context) mo.Result[mo.Option[*model.Entry]]
