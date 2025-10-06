package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ss49919201/keeput/app/analyzer/internal/appslog"
	"github.com/ss49919201/keeput/app/analyzer/internal/config"
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

func handleRequest(ctx context.Context) error {
	log.Print("Hello")
	return nil
}

func main() {
	lambda.Start(handleRequest)
}
