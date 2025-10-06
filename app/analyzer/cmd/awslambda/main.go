package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ss49919201/keeput/app/analyzer/internal/appctx"
	"github.com/ss49919201/keeput/app/analyzer/internal/appotel"
	"github.com/ss49919201/keeput/app/analyzer/internal/appslog"
	"github.com/ss49919201/keeput/app/analyzer/internal/config"
	"go.opentelemetry.io/otel"
)

const (
	traceName = "github.com/ss49919201/keeput/app/cmd/awslambda"
)

func init() {
	appslog.Init()
	if err := config.InitForLocal(); err != nil {
		slog.Error("failed to init env for local", slog.String("err", err.Error()))
		os.Exit(1)
	}
}

func handleRequest(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			appotel.RecordError(ctx, err)
		}
	}()

	ctx = appctx.SetNow(ctx, time.Now())

	shutdownTraceProvider, err := appotel.InitTraceProvider(ctx)
	if err != nil {
		slog.Error("failed to construct otel trace provider", slog.String("error", err.Error()))
	}
	defer func() {
		if err := shutdownTraceProvider(ctx); err != nil {
			slog.Warn("failed to shutdown trace provider", slog.String("error", err.Error()))
		}
	}()

	ctx, span := otel.Tracer(traceName).Start(ctx, "Lambda Entrypoint")
	defer span.End()

	log.Print("Hello")
	return nil
}

func main() {
	lambda.Start(handleRequest)
}
