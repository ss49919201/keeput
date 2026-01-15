package persister

import (
	"context"

	"github.com/ss49919201/keeput/app/analyzer/internal/model"
)

type PersistEntry = func(context.Context, *model.Entry) error
