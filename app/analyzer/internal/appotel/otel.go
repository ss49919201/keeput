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

type typeShutdownProvider = func(context.Context) error

var shutdownTraceProvider typeShutdownProvider
var errInitTraceProvider error
var initOtelProviderOnce sync.Once

func InitTraceProvider(ctx context.Context) (typeShutdownProvider, error) {
	initOtelProviderOnce.Do(func() {
		exporter, err := otlptracehttp.New(ctx,
			otlptracehttp.WithInsecure(),
		)
		if err != nil {
			errInitTraceProvider = err
			return
		}
		resource, err := resource.New(ctx)
		if err != nil {
			errInitTraceProvider = err
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
	return shutdownTraceProvider, errInitTraceProvider
}

// TODO
func InitMeterProvider(ctx context.Context) typeShutdownProvider

func RecordError(ctx context.Context, err error) {
	trace.SpanFromContext(ctx).RecordError(err, trace.WithStackTrace(true))
}
