package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/ss49919201/keeput/app/analyzer/internal/adapter/controller/cli"
	"github.com/ss49919201/keeput/app/analyzer/internal/appctx"
	"github.com/ss49919201/keeput/app/analyzer/internal/appotel"
	"github.com/ss49919201/keeput/app/analyzer/internal/config"
	"go.opentelemetry.io/otel"
)

const (
	traceName = "github.com/ss49919201/keeput/app/cmd/cli"
)

func initSlog() {
	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: lo.Switch[string, slog.Level](strings.ToUpper(config.LogLevel())).
					Case("DEBUG", slog.LevelDebug).
					Case("INFO", slog.LevelInfo).
					Case("WARN", slog.LevelWarn).
					Case("ERROR", slog.LevelError).
					Default(slog.LevelWarn),
			}),
		),
	)
}

func init() {
	initSlog()
}

func run(ctx context.Context) error {
	shutdownTraceProvider, err := appotel.InitTraceProvider(ctx)
	if err != nil {
		slog.Error("failed to construct otel trace provider", slog.String("error", err.Error()))
	}
	defer shutdownTraceProvider(ctx)

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
	slog.Info("success command")
}
