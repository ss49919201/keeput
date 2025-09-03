package usecase

import (
	"context"
	"time"

	"github.com/samber/mo"
	"github.com/ss49919201/fight-op/app/analyzer/internal/appctx"
	"github.com/ss49919201/fight-op/app/analyzer/internal/model"
	"github.com/ss49919201/fight-op/app/analyzer/internal/port/fetcher"
	"github.com/ss49919201/fight-op/app/analyzer/internal/port/usecase"
)

func NewAnalyze(fetchAllByDate fetcher.FetchLatest) usecase.Analyze {
	return func(ctx context.Context, in *usecase.AnalyzeInput) mo.Result[*usecase.AnalyzeOutput] {
		return analyze(ctx, in, fetchAllByDate)
	}
}

func analyze(ctx context.Context, in *usecase.AnalyzeInput, fetchLatest fetcher.FetchLatest) mo.Result[*usecase.AnalyzeOutput] {
	latestEntry, err := fetchLatest(ctx).Get()
	if err != nil {
		return mo.Err[*usecase.AnalyzeOutput](err)
	}

	return mo.Ok(&usecase.AnalyzeOutput{
		IsGoalAchieved: latestEntry.IsPresent() &&
			model.IsGoalAchieved(latestEntry.MustGet().PublishedAt, appctx.GetNowOr(ctx, time.Now()), in.Goal),
	})
}
