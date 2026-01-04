package registory

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/fetcher/hatena"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/fetcher/zenn"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/locker/cfworker"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/notifier/discord"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/persister/s3"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/printer/stdout"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/fetcher"
	usecaseport "github.com/ss49919201/keeput/app/analyzer/internal/port/usecase"
	usecaseadapter "github.com/ss49919201/keeput/app/analyzer/internal/usecase"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

func NewAnalyzeUsecase(ctx context.Context) (usecaseport.Analyze, error) {
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	otelaws.AppendMiddlewares(&awsConfig.APIOptions)

	return usecaseadapter.NewAnalyze(
		[]fetcher.FetchLatestEntry{
			hatena.NewFetchLatestEntry(),
			zenn.NewFetchLatestEntry(),
		},
		stdout.PrintAnalysisReport,
		discord.NewNotifyAnalysisReport(),
		cfworker.NewAcquire(),
		cfworker.NewRelease(),
		s3.NewPersistAnalysisReport(awsConfig),
	), nil
}
