package notifier

import (
	"context"

	"github.com/ss49919201/keeput/app/analyzer/internal/model"
)

type NotifyAnalysisReport = func(context.Context, *model.AnalysisReport) error
