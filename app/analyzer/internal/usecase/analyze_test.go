package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/samber/mo"
	"github.com/ss49919201/fight-op/app/analyzer/internal/appctx"
	"github.com/ss49919201/fight-op/app/analyzer/internal/model"
	"github.com/ss49919201/fight-op/app/analyzer/internal/port/fetcher"
	"github.com/ss49919201/fight-op/app/analyzer/internal/port/usecase"
	"github.com/stretchr/testify/assert"
)

func TestAnalyze(t *testing.T) {
	type args struct {
		fetchAllByDate func(t *testing.T) fetcher.FetchAllByDate

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
				fetchAllByDate: func(t *testing.T) fetcher.FetchAllByDate {
					return func(ctx context.Context, from, to time.Time) mo.Result[[]*model.Entry] {
						assert.Equal(t, time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC), from)
						assert.Equal(t, time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC), to)

						return mo.Ok([]*model.Entry{
							{
								Title:       "Go 言語の slice について",
								Body:        "Go 言語の slice は参照型です。気をつけましょう。",
								PublishedAt: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
							},
						})
					}
				},
				ctx: appctx.SetNow(context.Background(), time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)),
				input: &usecase.AnalyzeInput{
					Goal: usecase.RecentWeek,
				},
			},
			mo.Ok(&usecase.AnalyzeOutput{
				IsGoalAchieved: true,
			}),
		},
		{
			"return results of not achieving goal",
			args{
				fetchAllByDate: func(t *testing.T) fetcher.FetchAllByDate {
					return func(ctx context.Context, from, to time.Time) mo.Result[[]*model.Entry] {
						assert.Equal(t, time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC), from)
						assert.Equal(t, time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC), to)

						return mo.Ok([]*model.Entry{})
					}
				},
				ctx: appctx.SetNow(context.Background(), time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)),
				input: &usecase.AnalyzeInput{
					Goal: usecase.RecentWeek,
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
				tt.args.fetchAllByDate(t),
			)(
				tt.args.ctx, tt.args.input,
			)
			assert.Equal(t, tt.want, got)
		})
	}
}
