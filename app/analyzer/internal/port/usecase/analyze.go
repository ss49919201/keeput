package usecase

import (
	"context"

	"github.com/samber/mo"
	"github.com/ss49919201/keeput/app/analyzer/internal/model"
)

type AnalyzeInput struct {
	Goal model.GoalType
}
type AnalyzeOutput struct {
	IsGoalAchieved bool
}

type Analyze = func(context.Context, *AnalyzeInput) mo.Result[*AnalyzeOutput]
