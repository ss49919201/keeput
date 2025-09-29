package appotel

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var shutdownTraceProvider func(context.Context) error
var errInitOtelProvider error
var initOtelProviderOnce sync.Once

func InitTraceProvider(ctx context.Context) (func(context.Context) error, error) {
	initOtelProviderOnce.Do(func() {
		exporter, err := otlptracehttp.New(ctx)
		if err != nil {
			errInitOtelProvider = err
			return
		}
		resource, err := resource.New(ctx)
		if err != nil {
			errInitOtelProvider = err
			return
		}
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithSyncer(exporter),
			sdktrace.WithResource(resource),
		)
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.TraceContext{})
		shutdownTraceProvider = tp.Shutdown
	})
	return shutdownTraceProvider, errInitOtelProvider
}

func RecordError(ctx context.Context, err error) {
	trace.SpanFromContext(ctx).RecordError(err, trace.WithStackTrace(true))
}
