package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/samber/mo"
	"github.com/ss49919201/keeput/app/analyzer/internal/appctx"
	"github.com/ss49919201/keeput/app/analyzer/internal/model"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/fetcher"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/locker"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/notifier"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/persister"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyze(t *testing.T) {
	type args struct {
		NewLatestEntryFetchers   func(t *testing.T) []fetcher.FetchLatestEntry
		NewNotifyAnalysisReport  func(t *testing.T) notifier.NotifyAnalysisReport
		NewAcquireLock           func(t *testing.T) locker.Acquire
		NewReleaseLock           func(t *testing.T) locker.Release
		NewPersistAnalysisReport func(t *testing.T) persister.PersistAnalysisReport

		ctx   context.Context
		input *usecase.AnalyzeInput
	}
	tests := []struct {
		name string
		args args
		want mo.Result[*usecase.AnalyzeOutput]
	}{
		{
			"return results of achieving goal",
			args{
				NewLatestEntryFetchers: func(t *testing.T) []fetcher.FetchLatestEntry {
					return []fetcher.FetchLatestEntry{
						func(ctx context.Context) mo.Result[mo.Option[*model.Entry]] {
							return mo.Ok(mo.Some(&model.Entry{
								Title:       "Go 言語の slice について",
								Body:        "Go 言語の slice は参照型です。気をつけましょう。",
								PublishedAt: time.Date(2025, 1, 9, 10, 0, 0, 0, time.UTC),
							}))
						},
					}
				},
				NewNotifyAnalysisReport: func(t *testing.T) notifier.NotifyAnalysisReport {
					return func(ctx context.Context, report *model.AnalysisReport) error {
						assert.Equal(t, &model.AnalysisReport{
							IsGoalAchieved: true,
							LatestEntry: mo.Some(&model.Entry{
								Title:       "Go 言語の slice について",
								Body:        "Go 言語の slice は参照型です。気をつけましょう。",
								PublishedAt: time.Date(2025, 1, 9, 10, 0, 0, 0, time.UTC),
							}),
						}, report)
						return nil
					}
				},
				NewAcquireLock: func(t *testing.T) locker.Acquire {
					return func(ctx context.Context, lockID string) mo.Result[bool] {
						assert.Equal(t, "usecase:analyze:2025-01-10", lockID)
						return mo.Ok(true)
					}
				},
				NewReleaseLock: func(t *testing.T) locker.Release {
					return func(ctx context.Context, lockID string) error {
						assert.Equal(t, "usecase:analyze:2025-01-10", lockID)
						return nil
					}
				},
				NewPersistAnalysisReport: func(t *testing.T) persister.PersistAnalysisReport {
					return func(ctx context.Context, report *model.AnalysisReport) error {
						assert.Equal(t, &model.AnalysisReport{
							IsGoalAchieved: true,
							LatestEntry: mo.Some(&model.Entry{
								Title:       "Go 言語の slice について",
								Body:        "Go 言語の slice は参照型です。気をつけましょう。",
								PublishedAt: time.Date(2025, 1, 9, 10, 0, 0, 0, time.UTC),
							}),
						}, report)
						return nil
					}
				},
				ctx: appctx.SetNow(context.Background(), time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)),
				input: &usecase.AnalyzeInput{
					Goal: model.GoalTypeRecentWeek,
				},
			},
			mo.Ok(&usecase.AnalyzeOutput{
				IsGoalAchieved: true,
			}),
		},
		{
			"return results of not achieving goal",
			args{
				NewLatestEntryFetchers: func(t *testing.T) []fetcher.FetchLatestEntry {
					return []fetcher.FetchLatestEntry{
						func(ctx context.Context) mo.Result[mo.Option[*model.Entry]] {
							return mo.Ok(mo.None[*model.Entry]())
						},
					}
				},
				NewNotifyAnalysisReport: func(t *testing.T) notifier.NotifyAnalysisReport {
					return func(ctx context.Context, report *model.AnalysisReport) error {
						assert.Equal(t, &model.AnalysisReport{
							IsGoalAchieved: false,
							LatestEntry:    mo.None[*model.Entry](),
						}, report)
						return nil
					}
				},
				NewAcquireLock: func(t *testing.T) locker.Acquire {
					return func(ctx context.Context, lockID string) mo.Result[bool] {
						assert.Equal(t, "usecase:analyze:2025-01-10", lockID)
						return mo.Ok(true)
					}
				},
				NewReleaseLock: func(t *testing.T) locker.Release {
					return func(ctx context.Context, lockID string) error {
						assert.Equal(t, "usecase:analyze:2025-01-10", lockID)
						return nil
					}
				},
				NewPersistAnalysisReport: func(t *testing.T) persister.PersistAnalysisReport {
					return func(ctx context.Context, report *model.AnalysisReport) error {
						assert.Equal(t, &model.AnalysisReport{
							IsGoalAchieved: false,
							LatestEntry:    mo.None[*model.Entry](),
						}, report)
						return nil
					}
				},
				ctx: appctx.SetNow(context.Background(), time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)),
				input: &usecase.AnalyzeInput{
					Goal: model.GoalTypeRecentWeek,
				},
			},
			mo.Ok(&usecase.AnalyzeOutput{
				IsGoalAchieved: false,
			}),
		},
		{
			"return error when all latestEntryFetchers fail",
			args{
				NewLatestEntryFetchers: func(t *testing.T) []fetcher.FetchLatestEntry {
					return []fetcher.FetchLatestEntry{
						func(ctx context.Context) mo.Result[mo.Option[*model.Entry]] {
							return mo.Err[mo.Option[*model.Entry]](assert.AnError)
						},
						func(ctx context.Context) mo.Result[mo.Option[*model.Entry]] {
							return mo.Err[mo.Option[*model.Entry]](assert.AnError)
						},
					}
				},
				NewNotifyAnalysisReport: func(t *testing.T) notifier.NotifyAnalysisReport {
					return func(ctx context.Context, report *model.AnalysisReport) error {
						t.Error("should not be called")
						return nil
					}
				},
				NewAcquireLock: func(t *testing.T) locker.Acquire {
					return func(ctx context.Context, lockID string) mo.Result[bool] {
						return mo.Ok(true)
					}
				},
				NewReleaseLock: func(t *testing.T) locker.Release {
					return func(ctx context.Context, lockID string) error {
						return nil
					}
				},
				NewPersistAnalysisReport: func(t *testing.T) persister.PersistAnalysisReport {
					return func(ctx context.Context, report *model.AnalysisReport) error {
						t.Error("should not be called")
						return nil
					}
				},
				ctx: appctx.SetNow(context.Background(), time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)),
				input: &usecase.AnalyzeInput{
					Goal: model.GoalTypeRecentWeek,
				},
			},
			mo.Err[*usecase.AnalyzeOutput](assert.AnError),
		},
		{
			"continue processing when some latestEntryFetchers fail",
			args{
				NewLatestEntryFetchers: func(t *testing.T) []fetcher.FetchLatestEntry {
					return []fetcher.FetchLatestEntry{
						func(ctx context.Context) mo.Result[mo.Option[*model.Entry]] {
							return mo.Ok(mo.Some(&model.Entry{
								Title:       "Javaについて",
								Body:        "JavaはJVMで動作します。",
								PublishedAt: time.Date(2025, 1, 9, 10, 0, 0, 0, time.UTC),
							}))
						},
						func(ctx context.Context) mo.Result[mo.Option[*model.Entry]] {
							return mo.Err[mo.Option[*model.Entry]](assert.AnError)
						},
					}
				},
				NewNotifyAnalysisReport: func(t *testing.T) notifier.NotifyAnalysisReport {
					return func(ctx context.Context, report *model.AnalysisReport) error {
						assert.Equal(t, &model.AnalysisReport{
							IsGoalAchieved: true,
							LatestEntry: mo.Some(&model.Entry{
								Title:       "Javaについて",
								Body:        "JavaはJVMで動作します。",
								PublishedAt: time.Date(2025, 1, 9, 10, 0, 0, 0, time.UTC),
							}),
						}, report)
						return nil
					}
				},
				NewAcquireLock: func(t *testing.T) locker.Acquire {
					return func(ctx context.Context, lockID string) mo.Result[bool] {
						return mo.Ok(true)
					}
				},
				NewReleaseLock: func(t *testing.T) locker.Release {
					return func(ctx context.Context, lockID string) error {
						return nil
					}
				},
				NewPersistAnalysisReport: func(t *testing.T) persister.PersistAnalysisReport {
					return func(ctx context.Context, report *model.AnalysisReport) error {
						assert.Equal(t, &model.AnalysisReport{
							IsGoalAchieved: true,
							LatestEntry: mo.Some(&model.Entry{
								Title:       "Javaについて",
								Body:        "JavaはJVMで動作します。",
								PublishedAt: time.Date(2025, 1, 9, 10, 0, 0, 0, time.UTC),
							}),
						}, report)
						return nil
					}
				},
				ctx: appctx.SetNow(context.Background(), time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)),
				input: &usecase.AnalyzeInput{
					Goal: model.GoalTypeRecentWeek,
				},
			},
			mo.Ok(&usecase.AnalyzeOutput{
				IsGoalAchieved: true,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAnalyze(
				tt.args.NewLatestEntryFetchers(t),
				tt.args.NewNotifyAnalysisReport(t),
				tt.args.NewAcquireLock(t),
				tt.args.NewReleaseLock(t),
				tt.args.NewPersistAnalysisReport(t),
			)(
				tt.args.ctx, tt.args.input,
			)
			if tt.want.IsError() {
				require.True(t, got.IsError(), "expected error but got success")
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
