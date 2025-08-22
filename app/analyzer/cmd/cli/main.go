package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ss49919201/fight-op/app/analyzer/internal/adapter/controller/cli"
)

func init() {
	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelWarn, // TODO: load from config
			}),
		),
	)
}

func main() {
	err := cli.Analyze(context.Background())
	if err != nil {
		slog.Error("failed to run cli program", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("success command")
}
