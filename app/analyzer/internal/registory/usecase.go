package registory

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/samber/lo"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/fetcher/hatena"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/fetcher/zenn"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/locker/cfworker"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/notifier/discord"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/persister/s3"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/printer/stdout"
	"github.com/ss49919201/keeput/app/analyzer/internal/model"
	usecaseport "github.com/ss49919201/keeput/app/analyzer/internal/port/usecase"
	usecaseadapter "github.com/ss49919201/keeput/app/analyzer/internal/usecase"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

func NewAnalyzeUsecase(ctx context.Context) (usecaseport.Analyze, error) {
	var entryPlatformType model.EntryPlatformType
	for entryPlatform := range model.EntryPlatformIteratorOrderByPriorityAsc() {
		entryPlatformType = entryPlatform.Type()
		break
	}

	fetcher := lo.If(
		entryPlatformType == model.EntryPlatformTypeHatena, hatena.NewFetchLatestEntry(),
	).ElseIf(
		entryPlatformType == model.EntryPlatformTypeZenn, zenn.NewFetchLatestEntry(),
	).Else(hatena.NewFetchLatestEntry())

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	otelaws.AppendMiddlewares(&awsConfig.APIOptions)

	return usecaseadapter.NewAnalyze(
		fetcher,
		stdout.PrintAnalysisReport,
		discord.NewNotifyAnalysisReport(),
		cfworker.NewAcquire(),
		cfworker.NewRelease(),
		s3.NewPersistAnalysisReport(awsConfig),
	), nil
}
