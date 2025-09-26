package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/samber/mo"
	"github.com/ss49919201/keeput/app/analyzer/internal/appctx"
	"github.com/ss49919201/keeput/app/analyzer/internal/model"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/fetcher"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/locker"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/notifier"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/printer"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/usecase"
)

func NewAnalyze(fetchAllByDate fetcher.FetchLatest, printAnalysisReport printer.PrintAnalysisReport, notify notifier.NotifyAnalysisReport, acquireLock locker.Acquire, releaseLock locker.Release) usecase.Analyze {
	return func(ctx context.Context, in *usecase.AnalyzeInput) mo.Result[*usecase.AnalyzeOutput] {
		return analyze(ctx, in, fetchAllByDate, printAnalysisReport, notify, acquireLock, releaseLock)
	}
}

const lockIDPrefixAnalyze = "usecase:analyze"

func analyze(ctx context.Context, in *usecase.AnalyzeInput, fetchLatest fetcher.FetchLatest, printAnalysisReport printer.PrintAnalysisReport, notify notifier.NotifyAnalysisReport, acquireLock locker.Acquire, releaseLock locker.Release) mo.Result[*usecase.AnalyzeOutput] {
	lockID := lockIDPrefixAnalyze + ":" + appctx.GetNowOr(ctx, time.Now()).Format(time.DateOnly)
	acquireLockResult := acquireLock(ctx, lockID)
	if acquireLockResult.IsError() {
		return mo.Err[*usecase.AnalyzeOutput](fmt.Errorf("failed to lock: %w", acquireLockResult.Error()))
	}
	if !acquireLockResult.MustGet() {
		return mo.Err[*usecase.AnalyzeOutput](errors.New("lock already acquired"))
	}
	defer func() {
		if err := releaseLock(ctx, lockID); err != nil {
			slog.Warn("failed release lock")
		}
	}()

	latestEntry, err := fetchLatest(ctx).Get()
	if err != nil {
		return mo.Err[*usecase.AnalyzeOutput](fmt.Errorf("failed to fetch latest entry: %w", err))
	}

	report := model.Analyze(latestEntry, appctx.GetNowOr(ctx, time.Now()), in.Goal)

	if err := printAnalysisReport(report); err != nil {
		return mo.Err[*usecase.AnalyzeOutput](fmt.Errorf("failed to print anlysis report: %w", err))
	}

	if err := notify(ctx, report); err != nil {
		slog.Warn("failed to notify", "error", err)
	}

	return mo.Ok(&usecase.AnalyzeOutput{
		IsGoalAchieved: report.IsGoalAchieved,
	})
}
