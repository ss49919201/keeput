package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/ss49919201/keeput/app/analyzer/internal/appctx"
	"github.com/ss49919201/keeput/app/analyzer/internal/appotel"
	"github.com/ss49919201/keeput/app/analyzer/internal/appslog"
	"github.com/ss49919201/keeput/app/analyzer/internal/config"
	"github.com/ss49919201/keeput/app/analyzer/internal/model"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/usecase"
	"github.com/ss49919201/keeput/app/analyzer/internal/registory"
	"go.opentelemetry.io/otel"
)

const (
	traceName = "github.com/ss49919201/keeput/app/cmd/cli"
)

func init() {
	appslog.Init()
	if err := config.InitForLocal(); err != nil {
		slog.Error("failed to init env for local", slog.String("err", err.Error()))
		os.Exit(1)
	}
}

func run(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			appotel.RecordSpanError(ctx, err)
		}
	}()

	ctx = appctx.SetNow(context.Background(), time.Now())

	shutdownTraceProvider, err := appotel.InitTraceProvider(ctx)
	if err != nil {
		slog.Error("failed to construct otel trace provider", slog.String("error", err.Error()))
	}
	defer func() {
		if err := shutdownTraceProvider(ctx); err != nil {
			slog.Warn("failed to shutdown trace provider", slog.String("error", err.Error()))
		}
	}()
	shutdownMeterProvider, err := appotel.InitMeterProvider(ctx)
	if err != nil {
		slog.Error("failed to construct otel meter provider", slog.String("error", err.Error()))
	}
	defer func() {
		if err := shutdownMeterProvider(ctx); err != nil {
			slog.Warn("failed to shutdown meter provider", slog.String("error", err.Error()))
		}
	}()

	ctx, span := otel.Tracer(traceName).Start(ctx, "CLI Entrypoint")
	defer span.End()

	analyze, err := registory.NewAnalyzeUsecase(ctx)
	if err != nil {
		return err
	}
	result := analyze(ctx, &usecase.AnalyzeInput{
		Goal: model.GoalTypeRecentWeek,
	})
	if result.IsError() {
		return result.Error()
	}

	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		slog.Error("failed to run cli program", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
