package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
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

func initEnvForLocal() error {
	if !config.IsLocal() {
		return nil
	}
	return godotenv.Load()
}

func init() {
	initSlog()
	if err := initEnvForLocal(); err != nil {
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
