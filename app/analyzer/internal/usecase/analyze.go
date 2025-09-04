package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/samber/mo"
	"github.com/ss49919201/keeput/app/analyzer/internal/appctx"
	"github.com/ss49919201/keeput/app/analyzer/internal/model"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/fetcher"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/printer"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/usecase"
)

func NewAnalyze(fetchAllByDate fetcher.FetchLatest, printAnalysisReport printer.PrintAnalysisReport) usecase.Analyze {
	return func(ctx context.Context, in *usecase.AnalyzeInput) mo.Result[*usecase.AnalyzeOutput] {
		return analyze(ctx, in, fetchAllByDate, printAnalysisReport)
	}
}

func analyze(ctx context.Context, in *usecase.AnalyzeInput, fetchLatest fetcher.FetchLatest, printAnalysisReport printer.PrintAnalysisReport) mo.Result[*usecase.AnalyzeOutput] {
	latestEntry, err := fetchLatest(ctx).Get()
	if err != nil {
		return mo.Err[*usecase.AnalyzeOutput](fmt.Errorf("failed to fetch latest entry: %w", err))
	}

	report := model.Analyze(latestEntry, appctx.GetNowOr(ctx, time.Now()), in.Goal)

	var errs error
	if err := printAnalysisReport(report); err != nil {
		errs = errors.Join(errs, fmt.Errorf("failed to print anlysis report: %w", err))
	}

	if err != nil {
		return mo.Err[*usecase.AnalyzeOutput](fmt.Errorf("failed to fetch latest entry: %w", err))
	}

	return mo.Ok(&usecase.AnalyzeOutput{
		IsGoalAchieved: report.IsGoalAchieved,
	})
}
