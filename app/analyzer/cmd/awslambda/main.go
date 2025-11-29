package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/samber/lo"
	"github.com/ss49919201/keeput/app/analyzer/internal/appctx"
	"github.com/ss49919201/keeput/app/analyzer/internal/appotel"
	"github.com/ss49919201/keeput/app/analyzer/internal/appslog"
	"github.com/ss49919201/keeput/app/analyzer/internal/model"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/usecase"
	"github.com/ss49919201/keeput/app/analyzer/internal/registory"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
)

func init() {
	appslog.Init()
}

type goalType string

const (
	goalTypeRecentWeek  goalType = "recent_week"
	goalTypeRecentMonth goalType = "recent_month"
)

type payload struct {
	GoalType goalType
}

func parseGoalType(typ goalType) model.GoalType {
	return lo.Switch[goalType, model.GoalType](typ).
		Case(goalTypeRecentWeek, model.GoalTypeRecentWeek).
		Case(goalTypeRecentMonth, model.GoalTypeRecentMonth).
		Default(model.GoalTypeRecentWeek)
}

func handleRequest(ctx context.Context, payload payload) (err error) {
	defer func() {
		if err != nil {
			appotel.RecordError(ctx, err)
		}
	}()

	ctx = appctx.SetNow(ctx, time.Now())

	analyze, err := registory.NewAnalyzeUsecase(ctx)
	if err != nil {
		return err
	}
	result := analyze(ctx, &usecase.AnalyzeInput{
		Goal: parseGoalType(payload.GoalType),
	})
	if result.IsError() {
		return result.Error()
	}

	return nil
}

func main() {
	ctx := context.Background()
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
	lambda.Start(otellambda.InstrumentHandler(handleRequest, otellambda.WithPropagator(xray.Propagator{})))
}
