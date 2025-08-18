package main

import (
	"log/slog"
	"os"
)

func init() {
	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(os.Stdout, nil),
		),
	)
}

func main() {
	result := run()
	if result.IsError() {
		slog.Error("failed to run cli program", slog.String("error", result.Error().Error()))
		os.Exit(1)
	}
	slog.Info("success")
}
