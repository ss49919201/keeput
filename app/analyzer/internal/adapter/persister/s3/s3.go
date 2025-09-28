package s3

import (
	"context"

	"github.com/ss49919201/keeput/app/analyzer/internal/model"
)

func persist(ctx context.Context, report *model.AnalysisReport) error
