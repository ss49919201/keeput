package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/samber/mo"
	"github.com/samber/mo/result"
	"github.com/ss49919201/keeput/app/analyzer/internal/appctx"
	"github.com/ss49919201/keeput/app/analyzer/internal/model"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/fetcher"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/locker"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/notifier"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/persister"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

func NewAnalyze(latestEntryFetchers []fetcher.FetchLatestEntry, notifyAnalysisReport notifier.NotifyAnalysisReport, acquireLock locker.Acquire, releaseLock locker.Release, persistAnalysisReport persister.PersistAnalysisReport) usecase.Analyze {
	return func(ctx context.Context, in *usecase.AnalyzeInput) mo.Result[*usecase.AnalyzeOutput] {
		return analyze(ctx, in, latestEntryFetchers, notifyAnalysisReport, acquireLock, releaseLock, persistAnalysisReport)
	}
}

const (
	lockIDPrefixAnalyze = "usecase:analyze"
	meterName           = "github.com/ss49919201/keeput/app/analyzer/internal/usecase"
)

var (
	meter               = otel.Meter(meterName)
	counterGoalAchieved = sync.OnceValue(func() metric.Int64Counter {
		counter, err := meter.Int64Counter(
			"goal.achieved",
			metric.WithDescription("The number of goal achieved"),
		)
		if err != nil {
			slog.Error("failed to construct goal achieved counter", slog.String("error", err.Error()))
		}
		return counter
	})
)

func analyze(ctx context.Context, in *usecase.AnalyzeInput, latestEntryFetchers []fetcher.FetchLatestEntry, notifyAnalysisReport notifier.NotifyAnalysisReport, acquireLock locker.Acquire, releaseLock locker.Release, persistAnalysisReport persister.PersistAnalysisReport) mo.Result[*usecase.AnalyzeOutput] {
	// NOTE: defer でのロック解放遅延を analyze のブロックで行いたいため、ロック処理は result.Pipe に含めない
	lockID := lockIDPrefixAnalyze + ":" + appctx.GetNowOr(ctx, time.Now()).Format(time.DateOnly)
	acquired, err := acquireLock(ctx, lockID).Get()
	if err != nil {
		return mo.Err[*usecase.AnalyzeOutput](err)
	}
	if !acquired {
		return mo.Err[*usecase.AnalyzeOutput](errors.New("lock already acquired"))
	}
	defer func() {
		if err := releaseLock(ctx, lockID); err != nil {
			slog.Warn("failed release lock")
		}
	}()

	return result.Pipe5(
		func() mo.Result[mo.Option[*model.Entry]] {
			resultCh := make(chan mo.Result[mo.Option[*model.Entry]], len(latestEntryFetchers))
			for _, fetch := range latestEntryFetchers {
				go func() {
					resultCh <- fetch(ctx)
				}()
			}
			entries := make([]*model.Entry, 0, len(latestEntryFetchers))
			var errs []error
			for range len(latestEntryFetchers) {
				entry, err := (<-resultCh).Get()
				if err != nil {
					errs = append(errs, err)
					continue
				}
				if entry.IsNone() {
					continue
				}
				entries = append(entries, entry.MustGet())
			}
			if len(errs) == len(latestEntryFetchers) {
				return mo.Err[mo.Option[*model.Entry]](fmt.Errorf("all entry fetch operations failed: %w", errors.Join(errs...)))
			} else if len(errs) > 0 && len(errs) < len(latestEntryFetchers) {
				slog.Warn("some entry fetch operations failed", slog.String("error", errors.Join(errs...).Error()))
			}
			return mo.Ok(model.Latest(entries))
		}(),
		result.Map(func(entry mo.Option[*model.Entry]) *model.AnalysisReport {
			return model.Analyze(entry, appctx.GetNowOr(ctx, time.Now()), in.Goal)
		}),
		result.Map(func(report *model.AnalysisReport) *model.AnalysisReport {
			if err := persistAnalysisReport(ctx, report); err != nil {
				slog.Warn("failed to persist analysis report", slog.String("error", err.Error()))
			}
			return report
		}),
		result.Map(func(report *model.AnalysisReport) *model.AnalysisReport {
			if err := notifyAnalysisReport(ctx, report); err != nil {
				slog.Warn("failed to notify analysis report", "error", err)
			}
			return report
		}),
		result.Map(func(report *model.AnalysisReport) *model.AnalysisReport {
			if report.IsGoalAchieved {
				counterGoalAchieved().Add(ctx, 1)
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
