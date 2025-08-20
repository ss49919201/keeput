package usecase

import (
	"context"
	"time"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/ss49919201/fight-op/app/analyzer/internal/date"
	"github.com/ss49919201/fight-op/app/analyzer/internal/port/fetcher"
	"github.com/ss49919201/fight-op/app/analyzer/internal/port/usecase"
)

func NewAnalyze(fetchAllByDate fetcher.FetchAllByDate) usecase.Analyze {
	return func(ctx context.Context, in *usecase.AnalyzeInput) mo.Result[*usecase.AnalyzeOutput] {
		return analyze(ctx, in, fetchAllByDate)
	}
}

func analyze(ctx context.Context, in *usecase.AnalyzeInput, fetchAllByDate fetcher.FetchAllByDate) mo.Result[*usecase.AnalyzeOutput] {
	now := time.Now()
	beginningOfToday := date.BeginningOfDay(now)
	beginningOfBeforeXdays := date.AddDays(
		beginningOfToday,
		lo.If(in.Goal == usecase.RecentWeek, -7).
			ElseIf(in.Goal == usecase.RecentMonth, -30).
			Else(-7),
	)

	fetched := fetchAllByDate(ctx, beginningOfBeforeXdays, beginningOfToday)
	if fetched.IsError() {
		return mo.Err[*usecase.AnalyzeOutput](fetched.Error())
	}

	return mo.Ok(&usecase.AnalyzeOutput{
		IsGoalAchieved: len(fetched.MustGet()) > 0,
	})
}
