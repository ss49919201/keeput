package notifier

import (
    "context"

    "github.com/ss49919201/keeput/app/analyzer/internal/model"
)

// Notify sends a notification based on the analysis report.
// Implementations may ignore some fields or use just IsGoalAchieved.
type Notify = func(ctx context.Context, report *model.AnalysisReport) error

