package persister

import (
	"context"

	"github.com/ss49919201/keeput/app/analyzer/internal/model"
)

type PersistAnalysisReport = func(context.Context, *model.AnalysisReport) error
