package repository

import (
	"time"

	"github.com/samber/mo"
	"github.com/ss49919201/fight-op/app/analyzer/internal/model"
)

type FetchAllByDate = func(from, to time.Time) mo.Result[[]*model.Entry]
type FetchLatest = func() mo.Result[mo.Option[*model.Entry]]
