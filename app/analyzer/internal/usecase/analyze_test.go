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
		NewFetchLatest func(t *testing.T) fetcher.FetchLatest

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
							PublishedAt: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
						}))
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
			)(
				tt.args.ctx, tt.args.input,
			)
			assert.Equal(t, tt.want, got)
		})
	}
}
