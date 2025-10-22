package usecase

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/samber/mo"
	"github.com/samber/mo/result"
	"github.com/ss49919201/keeput/app/analyzer/internal/appctx"
	"github.com/ss49919201/keeput/app/analyzer/internal/model"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/fetcher"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/locker"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/notifier"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/persister"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/printer"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/usecase"
)

func NewAnalyze(fetchLatestEntry fetcher.FetchLatestEntry, printAnalysisReport printer.PrintAnalysisReport, notifyAnalysisReport notifier.NotifyAnalysisReport, acquireLock locker.Acquire, releaseLock locker.Release, persistAnalysisReport persister.PersistAnalysisReport) usecase.Analyze {
	return func(ctx context.Context, in *usecase.AnalyzeInput) mo.Result[*usecase.AnalyzeOutput] {
		return analyze(ctx, in, fetchLatestEntry, printAnalysisReport, notifyAnalysisReport, acquireLock, releaseLock, persistAnalysisReport)
	}
}

const lockIDPrefixAnalyze = "usecase:analyze"

func analyze(ctx context.Context, in *usecase.AnalyzeInput, fetchLatestEntry fetcher.FetchLatestEntry, printAnalysisReport printer.PrintAnalysisReport, notifyAnalysisReport notifier.NotifyAnalysisReport, acquireLock locker.Acquire, releaseLock locker.Release, persistAnalysisReport persister.PersistAnalysisReport) mo.Result[*usecase.AnalyzeOutput] {
	lockID := lockIDPrefixAnalyze + ":" + appctx.GetNowOr(ctx, time.Now()).Format(time.DateOnly)
	// NOTE: defer でのロック解放遅延を analyze のブロックで行いたいため、ロック処理は result.Pipe に含めない
	locked, err := acquireLock(ctx, lockID).Get()
	if err != nil {
		return mo.Err[*usecase.AnalyzeOutput](err)
	}
	if !locked {
		return mo.Err[*usecase.AnalyzeOutput](errors.New("lock already acquired"))
	}
	defer func() {
		if err := releaseLock(ctx, lockID); err != nil {
			slog.Warn("failed release lock")
		}
	}()

	return result.Pipe3(
		fetchLatestEntry(ctx),
		result.Map(func(entry mo.Option[*model.Entry]) *model.AnalysisReport {
			return model.Analyze(entry, appctx.GetNowOr(ctx, time.Now()), in.Goal)
		}),
		result.Map(func(report *model.AnalysisReport) *model.AnalysisReport {
			if err := persistAnalysisReport(ctx, report); err != nil {
				slog.Warn("failed to persist analysis report", slog.String("error", err.Error()))
			}
			if err := printAnalysisReport(report); err != nil {
				slog.Warn("failed to print anlysis report", slog.String("error", err.Error()))
			}
			if err := notifyAnalysisReport(ctx, report); err != nil {
				slog.Warn("failed to notify analysis report", "error", err)
			}
			return report
		}),
		result.Map(func(report *model.AnalysisReport) *usecase.AnalyzeOutput {
			return &usecase.AnalyzeOutput{
				IsGoalAchieved: report.IsGoalAchieved,
			}
		}),
	)
}
