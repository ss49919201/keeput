package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/samber/lo"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/fetcher/hatena"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/fetcher/zenn"
	"github.com/ss49919201/keeput/app/analyzer/internal/model"
	usecaseport "github.com/ss49919201/keeput/app/analyzer/internal/port/usecase"
	usecaseadapter "github.com/ss49919201/keeput/app/analyzer/internal/usecase"
)

type analyezeOutput struct {
	IsGoalAchieved bool `json:"is_goal_achieved"`
}

func Analyze(ctx context.Context) error {
	var entryPlatformType model.EntryPlatformType
	for entryPlatform := range model.EntryPlatformIteratorOrderByPriorityAsc() {
		entryPlatformType = entryPlatform.Type()
		break
	}

	fetcher := lo.If(
		entryPlatformType == model.EntryPlatformTypeHatena, hatena.NewFetchLatest(),
	).ElseIf(
		entryPlatformType == model.EntryPlatformTypeZenn, zenn.NewFetchLatest(),
	).Else(hatena.NewFetchLatest())

	result := usecaseadapter.NewAnalyze(
		fetcher,
	)(ctx, &usecaseport.AnalyzeInput{})
	if result.IsError() {
		return result.Error()
	}

	output := analyezeOutput{
		IsGoalAchieved: result.MustGet().IsGoalAchieved,
	}

	b, err := json.Marshal(output)
	if err != nil {
		return err
	}

	// TODO: use exporter
	fmt.Println(string(b))

	return nil
}
