package cli

import (
	"context"

	"github.com/samber/lo"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/fetcher/hatena"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/fetcher/zenn"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/locker/cfworker"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/printer/stdout"
	"github.com/ss49919201/keeput/app/analyzer/internal/model"
	usecaseport "github.com/ss49919201/keeput/app/analyzer/internal/port/usecase"
	usecaseadapter "github.com/ss49919201/keeput/app/analyzer/internal/usecase"
)

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
		stdout.PrintAnalysisReport,
		cfworker.NewAcquire(),
		cfworker.NewRelease(),
	)(ctx, &usecaseport.AnalyzeInput{})
	if result.IsError() {
		return result.Error()
	}

	return nil
}
