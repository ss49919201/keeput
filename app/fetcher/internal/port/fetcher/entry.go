package repository

import (
	"time"

	"github.com/samber/mo"
	"github.com/ss49919201/fight-op/app/fetcher/internal/model"
)

type FetchAllByDate = func(from, to time.Time) mo.Result[[]*model.EntryForRead]
type FetchLatest = func() mo.Result[mo.Option[*model.EntryForRead]]
