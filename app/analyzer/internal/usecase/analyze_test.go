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
	"github.com/ss49919201/keeput/app/analyzer/internal/port/printer"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/usecase"
	"github.com/stretchr/testify/assert"
)

func TestAnalyze(t *testing.T) {
type args struct {
        NewFetchLatest         func(t *testing.T) fetcher.FetchLatest
        NewPrintAnalysisReport func(t *testing.T) printer.PrintAnalysisReport
        NewNotify              func(t *testing.T) func(ctx context.Context, report *model.AnalysisReport) error
        NewAcquireLock         func(t *testing.T) locker.Acquire
        NewReleaseLock         func(t *testing.T) locker.Release

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
                NewFetchLatest: func(t *testing.T) fetcher.FetchLatest {
					return func(ctx context.Context) mo.Result[mo.Option[*model.Entry]] {
						return mo.Ok(mo.Some(&model.Entry{
							Title:       "Go 言語の slice について",
							Body:        "Go 言語の slice は参照型です。気をつけましょう。",
							PublishedAt: time.Date(2025, 1, 9, 10, 0, 0, 0, time.UTC),
						}))
					}
				},
				NewPrintAnalysisReport: func(t *testing.T) printer.PrintAnalysisReport {
					return func(report *model.AnalysisReport) error {
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
                NewNotify: func(t *testing.T) func(ctx context.Context, report *model.AnalysisReport) error {
                    return func(ctx context.Context, report *model.AnalysisReport) error { return nil }
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
                NewFetchLatest: func(t *testing.T) fetcher.FetchLatest {
					return func(ctx context.Context) mo.Result[mo.Option[*model.Entry]] {
						return mo.Ok(mo.None[*model.Entry]())
					}
				},
				NewPrintAnalysisReport: func(t *testing.T) printer.PrintAnalysisReport {
					return func(report *model.AnalysisReport) error {
						assert.Equal(t, &model.AnalysisReport{
							IsGoalAchieved: false,
							LatestEntry:    mo.None[*model.Entry](),
						}, report)

						return nil
					}
                },
                NewNotify: func(t *testing.T) func(ctx context.Context, report *model.AnalysisReport) error {
                    return func(ctx context.Context, report *model.AnalysisReport) error { return nil }
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
				ctx: appctx.SetNow(context.Background(), time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)),
				input: &usecase.AnalyzeInput{
					Goal: model.GoalTypeRecentWeek,
				},
			},
			mo.Ok(&usecase.AnalyzeOutput{
				IsGoalAchieved: false,
			}),
		},
	}
	for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := NewAnalyze(
                tt.args.NewFetchLatest(t),
                tt.args.NewPrintAnalysisReport(t),
                tt.args.NewNotify(t),
                tt.args.NewAcquireLock(t),
                tt.args.NewReleaseLock(t),
            )(
                tt.args.ctx, tt.args.input,
            )
			assert.Equal(t, tt.want, got)
		})
	}
}
