package usecase

import (
	"context"

	"github.com/samber/mo"
)

type GoalType int

const (
	RecentWeek  GoalType = iota + 1
	RecentMonth GoalType = iota + 1
)

type AnalyzeInput struct {
	Goal GoalType
}
type AnalyzeOutput struct {
	IsGoalAchieved bool

	// TODO: LatestEntryPublishedAt mo.Option[time.Time]
}

type Analyze = func(context.Context, *AnalyzeInput) mo.Result[*AnalyzeOutput]
