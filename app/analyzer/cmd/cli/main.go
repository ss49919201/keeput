package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/controller/cli"
	"github.com/ss49919201/keeput/app/analyzer/internal/appctx"
	"github.com/ss49919201/keeput/app/analyzer/internal/appotel"
	"github.com/ss49919201/keeput/app/analyzer/internal/appslog"
	"github.com/ss49919201/keeput/app/analyzer/internal/config"
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
			appotel.RecordError(ctx, err)
		}
	}()

	shutdownTraceProvider, err := appotel.InitTraceProvider(ctx)
	if err != nil {
		slog.Error("failed to construct otel trace provider", slog.String("error", err.Error()))
	}
	defer func() {
		if err := shutdownTraceProvider(ctx); err != nil {
			slog.Warn("failed to shutdown trace provider", slog.String("error", err.Error()))
		}
	}()

	ctx, span := otel.Tracer(traceName).Start(ctx, "CLI Entrypoint")
	defer span.End()
	if err := cli.Analyze(ctx); err != nil {
		return err
	}

	return nil
}

func main() {
	ctx := appctx.SetNow(context.Background(), time.Now())
	if err := run(ctx); err != nil {
		slog.Error("failed to run cli program", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
