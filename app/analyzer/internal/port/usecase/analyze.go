package usecase

import (
	"context"

	"github.com/samber/mo"
	"github.com/ss49919201/fight-op/app/analyzer/internal/model"
)

type AnalyzeInput struct {
	Goal model.GoalType
}
type AnalyzeOutput struct {
	IsGoalAchieved bool

	// TODO: LatestEntryPublishedAt mo.Option[time.Time]
}

type Analyze = func(context.Context, *AnalyzeInput) mo.Result[*AnalyzeOutput]
