package cli

import (
	"context"
	"fmt"

	"github.com/ss49919201/fight-op/app/analyzer/internal/adapter/fetcher/zenn"
	"github.com/ss49919201/fight-op/app/analyzer/internal/config"
	usecaseport "github.com/ss49919201/fight-op/app/analyzer/internal/port/usecase"
	usecaseadapter "github.com/ss49919201/fight-op/app/analyzer/internal/usecase"
)

func Analyze(ctx context.Context) error {
	result := usecaseadapter.NewAnalyze(
		zenn.NewFetchAllByDate(config.FeedURLZenn()),
	)(ctx, &usecaseport.AnalyzeInput{})
	if result.IsError() {
		return result.Error()
	}

	b, err := result.MarshalJSON()
	if err != nil {
		return err
	}

	// TODO: use exporter
	fmt.Println(string(b))

	return nil
}
